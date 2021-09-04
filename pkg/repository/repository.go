package repository

import (
	"github.com/jmoiron/sqlx"
	"qbot"
	"time"
)

type Bot interface {
	GetOrCreateSubscriber(chatId int64) (qbot.Subscriber, bool, error)
	SetSubscriberState(chatId int64, step string) error
	GetSubscriberState(chatId int64) (string, error)
	GetAyatByMailingDay(mailingDay int) (qbot.Ayat, error)
	GetActiveSubscribers() ([]qbot.Subscriber, error)
	DeactivateSubscribers([]int64) error
}

type Content interface {
	GetAyatByMailingDay(mailingDay int) (qbot.Ayat, error)
	GetAyatsBySuraNum(suraNum int) ([]qbot.Ayat, error)
	GetFavoriteAyats(chatId int64) ([]qbot.Ayat, error)
	GetAdjacentAyats(chatId int64, ayatId int) ([]qbot.Ayat, error)
	GetRandomPodcast() (qbot.Podcast, error)
	GetAyatById(chatId int64, ayatId int) (qbot.Ayat, error)
	AddToFavorite(chatId int64, ayatId int) error
	AyatIsFavorite(chatId int64, ayatId int) bool
	RemoveFromFavorite(chatId int64, ayatId int) error
}

type Prayer interface {
	GetPrayer(chatId int64, date time.Time) ([]qbot.Prayer, error)
	SubscriberHasCity(chatId int64) (bool, error)
}

type Repository struct {
	Bot
	Content
	Prayer
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Bot:     NewBotPostgres(db),
		Content: NewContentPostgres(db),
		Prayer:  NewPrayerPostgres(db),
	}
}
