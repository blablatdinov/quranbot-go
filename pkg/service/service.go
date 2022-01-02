package service

import (
	"qbot"
	"qbot/pkg/repository"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot interface {
	CreateSubscriber(chatId int64) (string, bool)
	SetSubscriberState(chatId int64, step string) error
	GetSubscriberState(chatId int64) (string, error)
	GetAyatByMailingDay(mailingDay int) (string, error)
	GetActiveSubscribers() ([]qbot.Subscriber, error)
	DeactivateSubscribers([]int64) error
	GetSubscribersCount(string) (int, error)
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
	GetMorningContentForTodayMailing() ([]qbot.MailingContent, error)
	UpdateDaysForSubscribers([]int64) error
}

type Prayer interface {
	GetPrayer(chatId int64, targetDate time.Time) (string, tgbotapi.InlineKeyboardMarkup, error)
	ChangePrayerStatus(prayerAtUserId int, status bool) (tgbotapi.InlineKeyboardMarkup, error)
	GetCityByName(cityName string) (qbot.City, error)
	ChangeCity(chatId int64, cityId int) error
	GetPrayersForMailing() ([]qbot.Answer, error)
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
