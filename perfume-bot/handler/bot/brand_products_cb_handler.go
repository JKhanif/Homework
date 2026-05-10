package bot_handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) BrandCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	data := update.CallbackQuery.Data

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.CallbackQuery.From.ID, Action: models.ChatActionTyping})

	parts := strings.Split(data, "_")
	var brandID string
	page := 0

	if strings.HasPrefix(data, "br_page_") {
		// parts: ["br", "page", "3", "1"]
		brandID = parts[2]
		page, _ = strconv.Atoi(parts[3])
	} else {
		// parts: ["brand", "3"]
		brandID = parts[1]
	}
	offset := page * pageSize

	brID, _ := strconv.Atoi(brandID)

	products, err := h.repo.GetProductsByBrandIDPage(ctx, brID, pageSize+1, offset)
	if err != nil {
		log.Printf("Error repo.GetProductsByBrandID: %v\n", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Ошибка, попробуйте позже.",
		})
		return
	}

	hasNext := len(products) > pageSize
	if hasNext {
		products = products[:pageSize]
	}

	for i, p := range products {
		var kb models.InlineKeyboardMarkup = models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Подробнее", CallbackData: fmt.Sprintf("detailed_%d", p.ID)},
					{Text: "В корзину", CallbackData: fmt.Sprintf("add_to_cart_%d", p.ID)},
				},
			},
		}

		brandTitle := p.Brand.Title
		if brandTitle == "" {
			brandTitle = "—"
		}

		if i == len(products)-1 {
			var navRow []models.InlineKeyboardButton
			if page > 0 {
				navRow = append(navRow, models.InlineKeyboardButton{Text: "←", CallbackData: fmt.Sprintf("br_page_%s_%d", brandID, page-1)})
			}
			if hasNext {
				navRow = append(navRow, models.InlineKeyboardButton{Text: "→", CallbackData: fmt.Sprintf("br_page_%s_%d", brandID, page+1)})
			}
			if len(navRow) > 0 {
				kb.InlineKeyboard = append(kb.InlineKeyboard, navRow)
			}
			kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
				{Text: "Главное меню", CallbackData: "main_menu"},
			})
		}

		_, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID: update.CallbackQuery.From.ID,
			Photo: &models.InputFileString{
				Data: *p.MainPhotoFailID,
			},
			Caption:     fmt.Sprintf("<b>%s</b> | %s\n\n%d₽", p.Title, brandTitle, p.Price),
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

func (h *Handler) BrandsCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.CallbackQuery.From.ID, Action: models.ChatActionTyping})

	brands, err := h.repo.GetAllBrands(ctx)
	if err != nil {
		log.Printf("Error repo.GetAllBrands: %v\n", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Ошибка, попробуйте позже.",
		})
		return
	}

	var kb models.InlineKeyboardMarkup
	kb.InlineKeyboard = make([][]models.InlineKeyboardButton, 0)
	for _, b := range brands {
		var row []models.InlineKeyboardButton
		row = append(row, models.InlineKeyboardButton{Text: b.Title, CallbackData: "brand_" + strconv.Itoa(b.ID)})
		kb.InlineKeyboard = append(kb.InlineKeyboard, row)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.From.ID,
		Text:        "Духи какого бренда вам интересны?",
		ReplyMarkup: kb,
	})
}
