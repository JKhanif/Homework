package http

import (
	"log"
	"net/http"
	"os"
	"perfume-bot/clients/minio"
	"perfume-bot/repository"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
)

type Handler struct {
	repo          *repository.Repository
	tgbot         *bot.Bot
	fileClient    *minio.Client
	chatID        int64
	minioPublicURL string
}

func NewHandler(repo *repository.Repository, tgbot *bot.Bot, fileClient *minio.Client, chatID int64, minioPublicURL string) *Handler {
	return &Handler{
		repo:           repo,
		tgbot:          tgbot,
		fileClient:     fileClient,
		chatID:         chatID,
		minioPublicURL: minioPublicURL,
	}
}

func (h *Handler) Run(port string) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	adminUser := getEnv("ADMIN_USER", "admin")
	adminPass := getEnv("ADMIN_PASSWORD", "admin")

	r.GET("/admin/*filepath", gin.BasicAuth(gin.Accounts{adminUser: adminPass}), func(c *gin.Context) {
		filepath := c.Param("filepath")
		if filepath == "" || filepath == "/" {
			filepath = "/index.html"
		}
		fullPath := "admin" + filepath
		data, err := os.ReadFile(fullPath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		ext := ""
		if i := strings.LastIndex(filepath, "."); i != -1 {
			ext = strings.ToLower(filepath[i:])
		}
		switch ext {
		case ".html":
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		case ".css":
			c.Data(http.StatusOK, "text/css; charset=utf-8", data)
		case ".js":
			c.Data(http.StatusOK, "application/javascript", data)
		case ".png":
			c.Data(http.StatusOK, "image/png", data)
		case ".jpg", ".jpeg":
			c.Data(http.StatusOK, "image/jpeg", data)
		case ".svg":
			c.Data(http.StatusOK, "image/svg+xml", data)
		default:
			c.Data(http.StatusOK, "text/plain; charset=utf-8", data)
		}
	})

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

	r.GET("/product/:id/photos", h.GetProductPhotosHandler)
	r.POST("/product/:id/photo-url", h.AddPhotoURLHandler)
	r.DELETE("/photo/:id", h.DeletePhotoHandler)
	r.PUT("/photo/:id/main/:product_id", h.SetMainPhotoHandler)

	r.POST("/upload", h.UploadHandler)
	r.GET("/uploads/:object", h.ServeFileHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
