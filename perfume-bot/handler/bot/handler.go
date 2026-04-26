package bot_handler

import (
	"context"
	"fmt"
	"log"
	"perfume-bot/repository"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/minio/minio-go/v7"
)

type Handler struct {
	repo  *repository.Repository
	minio *minio.Client
}

func NewHandler(repo *repository.Repository, minio *minio.Client) *Handler {
	return &Handler{
		repo:  repo,
		minio: minio,
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

	brands, err := h.repo.GetAllBrands(ctx)
	if err != nil {

	}

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
