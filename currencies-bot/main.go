package main

import (
	"context"
	"currencies/client/fxratesapi"
	"currencies/handler"
	"log"
	"os"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading secrets file")
	}

	fxRatesApiClient := fxratesapi.New(os.Getenv("FX_RATES_API_TOKEN="))
	handler := handler.New(fxRatesApiClient)

	ctx := context.Background()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Handle),
	}

	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}
