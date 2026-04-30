package http

import (
	"log"
	"net/http"
	"perfume-bot/model/api_model"
	"strconv"

	"github.com/gin-gonic/gin"
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

	// h.repo.SetProductCategories(c, id, )  // Если категории не пустые, добавить.

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

	c.JSON(http.StatusOK, products)
}
