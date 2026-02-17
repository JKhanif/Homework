package handlers

import (
	"context"
	"fmt"
	"log"
	"prayertimes/clients/aladhan"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var prayerNames = [5]string{"Фаджр", "Зухр", "Аср", "Магриб", "Иша"}

type Handler struct {
	aladhanClient *aladhan.Client
}

func New(aladhanClient *aladhan.Client) *Handler {
	return &Handler{
		aladhanClient: aladhanClient,
	}
}

func (h *Handler) PrayerNow(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	city := update.Message.Text

	res, err := h.aladhanClient.GetTodayPrayerTimesByCity(ctx, city)
	if err != nil {
		log.Println("error h.GetTodayPrayerTimesByCity: ", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка. Попробуйте позже.",
		})
		return
	}

	timeStr := []string{res.Data.Timings.Fajr, res.Data.Timings.Dhuhr, res.Data.Timings.Asr, res.Data.Timings.Maghrib, res.Data.Timings.Isha}

	nextPrayerFound := false
	resMessage := ""

	for i, v := range timeStr {
		now := time.Now()
		parsed, err := time.Parse("15:04", v)
		if err != nil {
			log.Println("error time.Parse: ", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка. Попробуйте позже.",
			})
			return
		}

		prayerTime := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, now.Location())

		if nextPrayerFound == false && now.Before(prayerTime) {
			nextPrayerFound = true
			resMessage += fmt.Sprintf("<b>%s: %s</b>\n", prayerNames[i], v)
			continue
		}

		resMessage += fmt.Sprintf("%s: %s\n", prayerNames[i], v)
	}

	now := time.Now().Format("02.01.2006")

	hijriTime := fmt.Sprintf("\n%s %s %s", res.Data.Date.Hijri.Day, res.Data.Date.Hijri.Month.En, res.Data.Date.Hijri.Year)
	resMessage = fmt.Sprintf("Город: %s\nРасписание на %s\n\n%s%s", city, now, resMessage, hijriTime)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      resMessage,
		ParseMode: models.ParseModeHTML,
	})
}
