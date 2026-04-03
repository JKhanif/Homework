package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"prayertimes/clients/aladhan"
	"prayertimes/handlers"
	"prayertimes/service"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if os.Getenv("TG_BOT_TOKEN") == "" {
		godotenv.Load()
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // "localhost:6379"
		Password: "",                      // если нет пароля
	})

	// проверка соединения
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Не удалось подключиться:", err)
	}
	fmt.Println("Подключено к Redis!")

	aladhanClient := aladhan.New(rdb)
	c := service.New(aladhanClient)
	h := handlers.New(c)

	opts := []bot.Option{
		bot.WithDefaultHandler(h.Handle),
	}

	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}
