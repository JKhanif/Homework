package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tasks-api/handlers"
	"tasks-api/repo"
	"tasks-api/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading secrets file")
		return
	}

	dbUrl := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	db, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal("No response from DB: ", err)
		return
	}

	repo := repo.NewRepo(db)
	service := service.New(repo)
	handler := handlers.New(service)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/tasks", handler.ReqTasksHandler)
	r.GET("/tasks/:id", handler.ReqTaskHandler)
	r.DELETE("/tasks/:id", handler.DeleteTaskHandler)
	r.PUT("/tasks/:id", handler.ChangeTaskHandler)
	r.POST("/tasks", handler.CreateTaskHandler)

	fmt.Println("Сервер запущен")
	r.Run(":8082")
}
