package http

import (
	"log"
	"net/http"
	"perfume-bot/model/api_model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategoryHandler(c *gin.Context) {
	var category api_model.CreateCategoryRequest

	err := c.ShouldBindJSON(&category)
	if err != nil {
		log.Printf("Error c.ShouldBindJSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON"})
		return
	}

	id, err := h.repo.CreateCategory(c, category)
	if err != nil {
		log.Printf("error repo.CreateCategory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении категории в БД"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status": "Категория добавлена!",
		"id":     id,
	})
}

func (h *Handler) UpdateCategoryHandler(c *gin.Context) { // cascade добавить
	var category api_model.UpdateCategoryRequest

	err := c.ShouldBindJSON(&category)
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

	err = h.repo.UpdateCategory(c, id, category)
	if err != nil {
		log.Printf("error repo.UpdateCategory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	c.JSON(http.StatusOK, "Изменено")
}

func (h *Handler) DeleteCategoryHandler(c *gin.Context) { // cascade добавить
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error strconv.Atoi(idStr): %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	err = h.repo.DeleteCategory(c, id)
	if err != nil {
		if err.Error() == "Нет в базе данных" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Не найдено"})
			return
		}

		log.Printf("error repo.DeleteCategory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	c.JSON(http.StatusOK, "Удалено")
}

func (h *Handler) GetCategoryHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error strconv.Atoi(idStr): %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
	}
	category, err := h.repo.GetCategoryByID(c, id)
	if err != nil {
		if err.Error() == "Не найдено" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Не найдено"})
			return
		}

		log.Printf("error repo.GetCategoryByID: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handler) GetAllCategoriesHandler(c *gin.Context) {
	categories, err := h.repo.GetAllCategories(c)
	if err != nil {
		log.Printf("error repo.GetAllCategories: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
