package bot_handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) CategoryCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	b.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: update.CallbackQuery.From.ID,
		Action: models.ChatActionTyping,
	})

	categoryID := strings.TrimPrefix(update.CallbackQuery.Data, "category_")
	products, err := h.repo.GetProductsByCategoryID(ctx, categoryID)
	if err != nil {
		log.Printf("Error repo.GetProductsByCategoryID: %v\n", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Ошибка, попробуйте позже.",
		})
		return
	}

	if len(products) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "В этой категории пока нет духов.",
		})
		return
	}

	for _, p := range products {
		var kb models.InlineKeyboardMarkup = models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Подробнее", CallbackData: fmt.Sprintf("detailed_%d", p.ID)},
				},
				{
					{Text: "В корзину", CallbackData: fmt.Sprintf("add_to_cart_%d", p.ID)},
				},
			},
		}

		_, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID: update.CallbackQuery.From.ID,
			Photo: &models.InputFileString{
				Data: p.MainPhotoFailID,
			},
			Caption:     fmt.Sprintf("<b>%s</b> | %s\n\n%d₽", p.Title, p.Brand.Title, p.Price),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		if err != nil {
			log.Printf("Error b.SendPhoto: %v\n", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Ошибка, попробуйте позже.",
			})
			return
		}
	}
}
