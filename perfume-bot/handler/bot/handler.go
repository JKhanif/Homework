package bot_handler

import (
	"context"
	"fmt"
	"log"
	"perfume-bot/repository"
	"strconv"
	"strings"

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

const pageSize = 3

func (h *Handler) CatalogCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	data := update.CallbackQuery.Data

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.CallbackQuery.From.ID, Action: models.ChatActionTyping})

	parts := strings.Split(data, "_")
	page, _ := strconv.Atoi(parts[1])
	offset := page * pageSize

	products, err := h.repo.GetAllProductsPage(ctx, pageSize+1, offset)
	if err != nil {
		log.Printf("Error repo.GetAllProductsPage: %v\n", err)
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
		if p.MainPhotoFailID == nil {
			continue
		}

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
				navRow = append(navRow, models.InlineKeyboardButton{Text: "←", CallbackData: fmt.Sprintf("catalog_%d", page-1)})
			}
			if hasNext {
				navRow = append(navRow, models.InlineKeyboardButton{Text: "→", CallbackData: fmt.Sprintf("catalog_%d", page+1)})
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

func (h *Handler) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Привет\n\nТут ты можешь купить парфюм",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Каталог", CallbackData: "catalog_0"},
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

func (h *Handler) MainMenuHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.From.ID,
		Text:   "Привет\n\nТут ты можешь купить парфюм",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Каталог", CallbackData: "catalog_0"},
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
