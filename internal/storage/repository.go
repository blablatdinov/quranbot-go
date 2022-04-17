package storage

import (
	"github.com/blablatdinov/quranbot-go/internal/core"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Bot interface {
	GetSubscriberByChatId(ChatId int64) (core.Subscriber, error)
	CreateSubscriber(ChatId int64, referralCode string) error
	ActivateSubscriber(chatId int64) error
}

type Content interface {
	GetAyatsBySuraNum(suraNum int) ([]core.Ayat, error)
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
