package http

import (
	"log"
	"net/http"
	"perfume-bot/clients/minio"
	"perfume-bot/repository"

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
	r.GET("/products", h.GetAllProductsHandler)
	r.GET("/product/:id", h.GetProductHandler)
	r.DELETE("/product/:id", h.DeleteProductHandler)
	r.PUT("/product/:id", h.ChangeProductHandler)
	r.POST("/products", h.CreateProductHandler)

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
