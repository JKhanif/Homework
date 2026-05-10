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

func (h *Handler) CategoryCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	data := update.CallbackQuery.Data

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.CallbackQuery.From.ID, Action: models.ChatActionTyping})

	parts := strings.Split(data, "_")
	var categoryID string
	page := 0

	if strings.HasPrefix(data, "cat_page_") {
		// parts: ["cat", "page", "5", "1"]
		categoryID = parts[2]
		page, _ = strconv.Atoi(parts[3])
	} else {
		// parts: ["category", "5"]
		categoryID = parts[1]
	}
	offset := page * pageSize

	catID, _ := strconv.Atoi(categoryID)

	products, err := h.repo.GetProductsByCategoryIDPage(ctx, catID, pageSize+1, offset)
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
				navRow = append(navRow, models.InlineKeyboardButton{Text: "←", CallbackData: fmt.Sprintf("cat_page_%s_%d", categoryID, page-1)})
			}
			if hasNext {
				navRow = append(navRow, models.InlineKeyboardButton{Text: "→", CallbackData: fmt.Sprintf("cat_page_%s_%d", categoryID, page+1)})
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

func (h *Handler) CategoriesCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.CallbackQuery.From.ID, Action: models.ChatActionTyping})

	categories, err := h.repo.GetAllCategories(ctx)
	if err != nil {
		log.Printf("Error repo.GetAllCategories: %v\n", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Ошибка, попробуйте позже.",
		})
		return
	}

	var kb models.InlineKeyboardMarkup
	kb.InlineKeyboard = make([][]models.InlineKeyboardButton, 0)
	for _, c := range categories {
		var row []models.InlineKeyboardButton
		row = append(row, models.InlineKeyboardButton{Text: c.Title, CallbackData: "category_" + strconv.Itoa(c.ID)})
		kb.InlineKeyboard = append(kb.InlineKeyboard, row)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.From.ID,
		Text:        "Выбирай по душе",
		ReplyMarkup: kb,
	})
}
