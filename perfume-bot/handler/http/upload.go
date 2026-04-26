package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdloadHandler(c *gin.Context) {
	from, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := from.File["files"]
	links := make([]string, 0, len(files))
	for _, file := range files {
		link, err := h.fileClient.UploadPhoto(c, file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		links = append(links, link)
	}

	c.JSON(http.StatusOK, gin.H{"files": links})
}
