package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	minio_cl "perfume-bot/clients/minio"
	bot_handler "perfume-bot/handler/bot"
	http_handler "perfume-bot/handler/http"
	"perfume-bot/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := pgxpool.New(ctx, os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatalf("Error pgxpool new: %v", err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Error postgre ping: %v", err)
	}

	mcl, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds: credentials.NewStaticV4(
			os.Getenv("MINIO_ACCESS_KEY"),
			os.Getenv("MINIO_SECRET_KEY"),
			"",
		),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("minio init error: %v", err)
	}

	repo := repository.New(pool)
	botHandler := bot_handler.NewHandler(repo, mcl)

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.DefaultHandler),
	}

	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"), opts...)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, botHandler.StartHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "catalog", bot.MatchTypeExact, botHandler.CatalogCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "categories", bot.MatchTypeExact, botHandler.CategoriesCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "category_", bot.MatchTypePrefix, botHandler.CategoryCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "brands", bot.MatchTypeExact, botHandler.BrandsCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "brand_", bot.MatchTypePrefix, botHandler.BrandCallbackHandler)

	setCommands(b, ctx)
	go b.Start(ctx)

	fileClient := minio_cl.New(mcl, os.Getenv("MINIO_BUCKET"), os.Getenv("MINIO_ENDPOINT"))

	chatID, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Invalid CHAT_ID: %v", err)
	}

	httpHandler := http_handler.NewHandler(repo, b, fileClient, chatID, os.Getenv("MINIO_PUBLIC_URL"))

	go func() {
		httpHandler.Run(os.Getenv("PORT"))
	}()

	fmt.Println("Started")

	// ждём Ctrl+C
	<-ctx.Done()

	fmt.Println("Shutting down...")
}

func setCommands(b *bot.Bot, ctx context.Context) {
	commands := []models.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		//{Command: "help", Description: "Помощь"},
		//{Command: "home", Description: "Главное меню"},
	}

	_, err := b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		log.Fatal(err)
	}
}
