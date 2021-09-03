package repository

import (
	"github.com/jmoiron/sqlx"
	"qbot"
)

type Bot interface {
	GetOrCreateSubscriber(chatId int64) (qbot.Subscriber, bool, error)
	SetSubscriberState(chatId int64, step string) error
	GetSubscriberState(chatId int64) (string, error)
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

type Repository struct {
	Bot
	Content
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Bot:     NewBotPostgres(db),
		Content: NewContentPostgres(db),
	}
}
