package main

import (
	"context"
	"os"
	"os/signal"
	"prayertimes/clients/aladhan"
	"prayertimes/handlers"
	"prayertimes/service"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	aladhanClient := aladhan.New()
	c := service.New(aladhanClient)
	h := handlers.New(c)

	opts := []bot.Option{
		bot.WithDefaultHandler(h.Handle),
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}
