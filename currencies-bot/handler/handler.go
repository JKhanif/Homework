package handler

import (
	"context"
	"currencies/client/fxratesapi"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	CurrencyRub = "RUB"
	CurrencyUsd = "USD"
	CurrencyTry = "TRY"
	CurrencySar = "SAR"
	CurrencyEur = "EUR"
)

func New(fraClient *fxratesapi.FraClient) *Handler {
	return &Handler{
		FraClient: fraClient,
	}
}

type Handler struct {
	FraClient *fxratesapi.FraClient
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.Message.Text == "/start" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Привет! Здесь ты можешь узнать курс валют, а также конвертировать их.",
		})
		return
	}

	h.HandleCurrencyConvert(ctx, b, update)
}

func (h *Handler) HandleCurrencyConvert(ctx context.Context, b *bot.Bot, update *models.Update) {
	currencyFrom := extractCurrency(update.Message.Text)

	if currencyFrom == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Поддерживаемые валюты: RUB, USD, EUR, TRY, SAR.",
		})
	}

	re := regexp.MustCompile(`\d+(\.\d+)?`)
	valString := re.FindString(update.Message.Text)

	val, err := strconv.Atoi(valString)
	if err != nil {
		fmt.Println("error strconv: ", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error, try again later.",
		})
		return
	}

	rates, err := h.FraClient.GetCurrencyRate(ctx, currencyFrom)
	if err != nil {
		fmt.Println("error FraClient.GetCurrencyRate: ", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error, try again later.",
		})
		return
	}

	var resp string = "<code>"
	var supportedCurrencies = []string{CurrencyRub, CurrencyUsd, CurrencyEur, CurrencyTry, CurrencySar}
	for _, cur := range supportedCurrencies {
		if currencyFrom != cur {
			resp += fmt.Sprintf("%s: %.2f\n", cur, rates[cur]*float32(val))
		}
	}

	resp += "</code>"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      resp,
		ParseMode: models.ParseModeHTML,
	})
}

func extractCurrency(s string) string {
	if strings.Contains(strings.ToLower(s), "руб") ||
		strings.Contains(strings.ToLower(s), "rub") || strings.Contains(s, "₽") {
		return CurrencyRub
	} else if strings.Contains(strings.ToLower(s), "usd") ||
		strings.Contains(s, "$") {
		return CurrencyUsd
	} else if strings.Contains(strings.ToLower(s), "try") ||
		strings.Contains(s, "₺") || strings.Contains(strings.ToLower(s), "лир") {
		return CurrencyTry
	} else if strings.Contains(strings.ToLower(s), "sar") ||
		strings.Contains(strings.ToLower(s), "риял") || strings.Contains(strings.ToLower(s), "riyal") {
		return CurrencySar
	} else if strings.Contains(strings.ToLower(s), "eur") ||
		strings.Contains(s, "€") || strings.Contains(strings.ToLower(s), "евро") {
		return CurrencyEur
	}

	return ""
}
