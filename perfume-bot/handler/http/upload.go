package http

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) ServeFileHandler(c *gin.Context) {
	objectName := c.Param("object")

	obj, err := h.fileClient.GetObject(c.Request.Context(), objectName)
	if err != nil {
		log.Printf("Error GetObject: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer obj.Close()

	objInfo, err := obj.Stat()
	if err != nil {
		log.Printf("Error object Stat: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.DataFromReader(http.StatusOK, objInfo.Size, objInfo.ContentType, obj, nil)
}

func (h *Handler) UploadHandler(c *gin.Context) {
	productIDStr := c.PostForm("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	firstIsMain := c.PostForm("is_main") != "false"

	if firstIsMain {
		err := h.repo.UnsetMainPhoto(c, productID)
		if err != nil {
			log.Printf("Error UnsetMainPhoto: %v", err)
		}
	}

	results := make([]gin.H, 0, len(files))

	for i, file := range files {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data, err := io.ReadAll(src)
		src.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		url, err := h.fileClient.UploadFromReader(c, bytes.NewReader(data), file.Filename, file.Size, file.Header.Get("Content-Type"))
		if err != nil {
			log.Printf("Error UploadFromReader: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "MinIO upload failed"})
			return
		}

		msg, err := h.tgbot.SendPhoto(c.Request.Context(), &bot.SendPhotoParams{
			ChatID: h.chatID,
			Photo: &models.InputFileUpload{
				Filename: file.Filename,
				Data:     bytes.NewReader(data),
			},
		})
		if err != nil {
			log.Printf("Error SendPhoto to Telegram: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Telegram send failed"})
			return
		}

		fileID := msg.Photo[0].FileID

		thisIsMain := i == 0 && firstIsMain

		photoID, err := h.repo.CreateProductPhoto(c, productID, url, fileID, thisIsMain)
		if err != nil {
			log.Printf("Error CreateProductPhoto: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB save failed"})
			return
		}

		results = append(results, gin.H{
			"photo_id":   photoID,
			"url":        url,
			"tg_file_id": fileID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"files": results})
}
