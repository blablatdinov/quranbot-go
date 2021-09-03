package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"qbot"
	"qbot/pkg/repository"
)

type Bot interface {
	CreateSubscriber(chatId int64) (string, bool)
	GetDefaultKeyboard() tgbotapi.ReplyKeyboardMarkup
	SetSubscriberState(chatId int64, step string) error
	GetSubscriberState(chatId int64) (string, error)
}

type Content interface {
	GetAyatByMailingDay(mailingDay int) (string, error)
	GetAyatBySuraAyatNum(chatId int64, query string, state string) (string, tgbotapi.InlineKeyboardMarkup, error)
	GetAyatById(chatId int64, ayatId int, state string) (string, tgbotapi.InlineKeyboardMarkup, error)
	GetFavoriteAyats(chatId int64) (string, tgbotapi.InlineKeyboardMarkup, error)
	GetFavoriteAyatsFromKeyboard(chatId int64, ayatId int) (string, tgbotapi.InlineKeyboardMarkup, error)
	GetRandomPodcast() (qbot.Podcast, error)
	AddToFavorite(chatId int64, ayatId int, state string) (string, tgbotapi.InlineKeyboardMarkup, error)
	RemoveFromFavorite(chatId int64, ayatId int, state string) (string, tgbotapi.InlineKeyboardMarkup, error)
}

type Prayer interface {
	GetPrayer(chatId int64) (string, error)
}

type Service struct {
	Bot
	Content
	Prayer
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Bot:     NewBotService(repos.Bot),
		Content: NewContentService(repos.Content),
		Prayer:  NewPrayerService(repos.Prayer),
	}
}