package handler

import (
	"context"
	"fmt"
	"log"
	"perfume-bot/repository"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) BrandsCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	b.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: update.CallbackQuery.From.ID,
		Action: models.ChatActionTyping,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:   "1, 2, 3",
	})
}

func (h *Handler) CategoryCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	answer := ""

	for _, p := range products {
		answer += fmt.Sprintf(
			"💎 %s\n🏷 Бренд: %s\n💰 Цена: %d\n\n%s\n\n",
			p.Title,
			p.Brand.Title,
			p.PriceKopecks,
			p.Description,
		)
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	b.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: update.CallbackQuery.From.ID,
		Action: models.ChatActionTyping,
	})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.From.ID,
		Text:      answer,
		ParseMode: models.ParseModeHTML,
	})
}

func (h *Handler) CategoriesCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	b.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: update.CallbackQuery.From.ID,
		Action: models.ChatActionTyping,
	})

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

func (h *Handler) CatalogCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer := ""

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
	b.SendChatAction(ctx, &bot.SendChatActionParams{
		ChatID: update.CallbackQuery.From.ID,
		Action: models.ChatActionTyping,
	})

	products, err := h.repo.GetAllProducts(ctx)
	if err != nil {
		log.Printf("Error repo.GetAllProducts: %v\n", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Ошибка, попробуйте позже.",
		})
		return
	}

	for _, p := range products {
		answer += fmt.Sprintf("<b>%s</b> | %s\n\n%d₽\n\n", p.Title, p.Brand.Title, p.PriceKopecks)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.From.ID,
		Text:      answer,
		ParseMode: models.ParseModeHTML,
	})
}

func (h *Handler) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Привет\n\nТут ты можешь купить парфюм",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Каталог", CallbackData: "catalog"},
					{Text: "Категории", CallbackData: "categories"},
				},
				{
					{Text: "Бренды", CallbackData: "brands"},
				},
				{
					{Text: "Корзина", CallbackData: "cart"},
				},
			},
		},
	})
	if err != nil {
		log.Printf("error b.SendMessage in StartHandler: %v\n", err)
	}
}

func (h *Handler) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Чтобы вызвать меню, нажмите /start",
	})
}
