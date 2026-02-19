package handlers

import (
	"context"
	"log"
	"prayertimes/service"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) PrayerNow(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	city := update.Message.Text
	res, err := h.service.ResponsePrayerTime(ctx, city)
	if err != nil {
		log.Println("error h.service.ResponsePrayerTime: ", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка. Попробуйте позже.",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      res,
		ParseMode: models.ParseModeHTML,
	})
}
