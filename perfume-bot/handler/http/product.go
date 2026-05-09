package http

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"perfume-bot/model/api_model"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) CreateProductHandler(c *gin.Context) {
	var product api_model.CreateProductRequest

	err := c.ShouldBindJSON(&product)
	if err != nil {
		log.Printf("Error c.ShouldBindJSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON"})
		return
	}

	id, err := h.repo.CreateProduct(c, product)
	if err != nil {
		log.Printf("error repo.CreateProduct: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении товара в БД"})
		return
	}

	err = h.repo.SetProductCategories(c, id, product.CategoryIDs)
	if err != nil {
		log.Printf("error repo.SetProductCategories: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении категорий товара"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status": "Товар добавлен!",
		"id":     id,
	})
}

func (h *Handler) UpdateProductHandler(c *gin.Context) {
	var product api_model.UpdateProductRequest

	err := c.ShouldBindJSON(&product)
	if err != nil {
		log.Printf("Error c.ShouldBindJSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error strconv.Atoi(idStr): %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	err = h.repo.UpdateProduct(c, id, product)
	if err != nil {
		log.Printf("error repo.UpdateProduct: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	if product.CategoryIDs != nil {
		err = h.repo.SetProductCategories(c, id, product.CategoryIDs)
		if err != nil {
			log.Printf("error repo.SetProductCategories: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, "Изменено")
}

func (h *Handler) DeleteProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error strconv.Atoi(idStr): %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	// 1. удаляем фото из MinIO
	photos, err := h.repo.GetProductPhotos(c, id)
	if err != nil {
		log.Printf("error GetProductPhotos: %v", err)
	}
	for _, ph := range photos {
		if h.fileClient != nil && ph.URL != "" && h.isMinioURL(ph.URL) {
			key := h.extractMinioKey(ph.URL)
			if key != "" {
				err := h.fileClient.DeleteObject(c, key)
				if err != nil {
					log.Printf("error DeleteObject from MinIO: %v", err)
				}
			}
		}
	}

	// 2. удаляем фото из БД
	err = h.repo.DeleteProductPhotos(c, id)
	if err != nil {
		log.Printf("error DeleteProductPhotos: %v", err)
	}

	// 3. удаляем связи с категориями
	err = h.repo.SetProductCategories(c, id, []int{})
	if err != nil {
		log.Printf("error SetProductCategories: %v", err)
	}

	// 4. удаляем сам продукт
	err = h.repo.DeleteProduct(c, id)
	if err != nil {
		if err.Error() == "Нет в базе данных" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Не найдено"})
			return
		}

		log.Printf("error repo.DeleteProduct: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	c.JSON(http.StatusOK, "Удалено")
}

func (h *Handler) GetProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error strconv.Atoi(idStr): %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	product, err := h.repo.GetProductByID(c, id)
	if err != nil {
		if err.Error() == "Не найдено" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Не найдено"})
			return
		}

		log.Printf("error repo.GetProductByID: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) GetAllProductsHandler(c *gin.Context) {
	products, err := h.repo.GetAllProducts(c)
	if err != nil {
		log.Printf("error repo.GetAllProducts: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	for i := range products {
		if products[i].MainPhotoURL != nil && (strings.HasPrefix(*products[i].MainPhotoURL, "http://localhost:9000/") || strings.HasPrefix(*products[i].MainPhotoURL, "http://minio:9000/")) {
			parts := strings.Split(*products[i].MainPhotoURL, "/")
			proxied := "/uploads/" + parts[len(parts)-1]
			products[i].MainPhotoURL = &proxied
		}
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) GetProductPhotosHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	photos, err := h.repo.GetProductPhotos(c, id)
	if err != nil {
		log.Printf("error repo.GetProductPhotos: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	// Проксируем URL через Go-сервер, т.к. MinIO не разрешает анонимный доступ
	for i := range photos {
		if strings.HasPrefix(photos[i].URL, "http://localhost:9000/") || strings.HasPrefix(photos[i].URL, "http://minio:9000/") {
			parts := strings.Split(photos[i].URL, "/")
			photos[i].URL = "/uploads/" + parts[len(parts)-1]
		}
	}

	c.JSON(http.StatusOK, photos)
}

func (h *Handler) DeletePhotoHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	// получить URL фото до удаления
	photo, err := h.repo.GetPhotoByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Не найдено"})
		return
	}

	// если фото в MinIO — удалить объект
	if h.fileClient != nil && photo.URL != "" {
		// URL формата "http://<endpoint>/<bucket>/<key>"
		// вытаскиваем key — часть после bucket/
		if h.isMinioURL(photo.URL) {
			key := h.extractMinioKey(photo.URL)
			if key != "" {
				err := h.fileClient.DeleteObject(c, key)
				if err != nil {
					log.Printf("error DeleteObject from MinIO: %v", err)
				}
			}
		}
	}

	err = h.repo.DeletePhoto(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Удалено")
}

func (h *Handler) isMinioURL(url string) bool {
	return h.minioPublicURL != "" && strings.HasPrefix(url, h.minioPublicURL)
}

func (h *Handler) extractMinioKey(url string) string {
	// URL формата "http://localhost:9000/images/<key>" или "http://minio:9000/images/<key>"
	// ищем второе вхождение "/" после "http://" чтобы пропустить хост, затем bucket
	parts := strings.SplitN(url, "/", 5)
	if len(parts) < 5 {
		return ""
	}
	// parts = ["http:", "", "localhost:9000", "images", "filename.jpg"]
	return parts[4]
}

func (h *Handler) SetMainPhotoHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный product_id"})
		return
	}

	err = h.repo.SetMainPhoto(c, id, productID)
	if err != nil {
		log.Printf("error repo.SetMainPhoto: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	c.JSON(http.StatusOK, "Главное фото обновлено")
}

type AddPhotoURLRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *Handler) AddPhotoURLHandler(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	var req AddPhotoURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url required"})
		return
	}

	// скачиваем изображение по URL
	resp, err := http.Get(req.URL)
	if err != nil {
		log.Printf("error downloading URL: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось скачать изображение"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL вернул HTTP " + strconv.Itoa(resp.StatusCode)})
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
		return
	}

	// определяем имя файла из URL
	filename := extractFilename(req.URL)
	contentType := resp.Header.Get("Content-Type")

	// загружаем в MinIO
	url, err := h.fileClient.UploadFromReader(c, bytes.NewReader(data), filename, int64(len(data)), contentType)
	if err != nil {
		log.Printf("error UploadFromReader: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MinIO upload failed"})
		return
	}

	// если у товара нет фото — делаем это главным
	existing, err := h.repo.GetProductPhotos(c, productID)
	if err != nil {
		log.Printf("error GetProductPhotos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}
	isMain := len(existing) == 0

	if isMain {
		h.repo.UnsetMainPhoto(c, productID)
	}

	// отправляем в Telegram для получения tg_file_id
	var tgFileID string
	if h.tgbot != nil {
		msg, err := h.tgbot.SendPhoto(c.Request.Context(), &bot.SendPhotoParams{
			ChatID: h.chatID,
			Photo: &models.InputFileUpload{
				Filename: filename,
				Data:     bytes.NewReader(data),
			},
		})
		if err == nil && len(msg.Photo) > 0 {
			tgFileID = msg.Photo[0].FileID
		} else {
			log.Printf("SendPhoto for URL photo failed (non-critical): %v", err)
		}
	}

	photoID, err := h.repo.CreateProductPhoto(c, productID, url, tgFileID, isMain)
	if err != nil {
		log.Printf("error CreateProductPhoto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo_id":   photoID,
		"url":        url,
		"tg_file_id": tgFileID,
	})
}

// извлекает имя файла из URL, либо генерирует image.jpg
func extractFilename(rawURL string) string {
	parts := strings.Split(rawURL, "/")
	last := parts[len(parts)-1]
	if strings.Contains(last, ".") && !strings.Contains(last, "?") {
		return last
	}
	return "image.jpg"
}
