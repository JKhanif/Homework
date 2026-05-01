package http

import (
	"log"
	"net/http"
	"perfume-bot/clients/minio"
	"perfume-bot/repository"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
)

type Handler struct {
	repo       *repository.Repository
	tgbot      *bot.Bot
	fileClient *minio.Client
}

func NewHandler(repo *repository.Repository, tgbot *bot.Bot, fileClient *minio.Client) *Handler {
	return &Handler{
		repo:       repo,
		tgbot:      tgbot,
		fileClient: fileClient,
	}
}

func (h *Handler) Run(port string) {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/products", h.GetAllProductsHandler)
	r.GET("/product/:id", h.GetProductHandler)
	r.DELETE("/product/:id", h.DeleteProductHandler)
	r.PUT("/product/:id", h.UpdateProductHandler)
	r.POST("/products", h.CreateProductHandler)

	r.GET("/brands", h.GetAllBrandsHandler)
	r.GET("/brand/:id", h.GetBrandHandler)
	r.DELETE("/brand/:id", h.DeleteBrandHandler)
	r.PUT("/brand/:id", h.UpdateBrandHandler)
	r.POST("/brands", h.CreateBrandHandler)

	r.GET("/categories", h.GetAllCategoriesHandler)
	r.GET("/category/:id", h.GetCategoryHandler)
	r.DELETE("/category/:id", h.DeleteCategoryHandler)
	r.PUT("/category/:id", h.UpdateCategoryHandler)
	r.POST("/categories", h.CreateCategoryHandler)

	r.POST("/upload", h.UpdloadHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
