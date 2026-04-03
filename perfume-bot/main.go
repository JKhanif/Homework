package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"perfume-bot/handler"
	"perfume-bot/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pool, err := pgxpool.New(ctx, os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatalf("Error pgxpool new: %v", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Error postgre ping: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	handler := handler.New(repo)

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.DefaultHandler),
	}

	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"), opts...)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, handler.StartHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "catalog", bot.MatchTypeExact, handler.CatalogCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "categories", bot.MatchTypeExact, handler.CategoriesCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "category_", bot.MatchTypePrefix, handler.CategoryCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "brands", bot.MatchTypeExact, handler.BrandsCallbackHandler)

	setCommands(b, ctx)
	b.Start(ctx)
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
