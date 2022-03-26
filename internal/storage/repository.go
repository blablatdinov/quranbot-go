package storage

import (
	"github.com/blablatdinov/quranbot-go/internal/core"

	"github.com/jmoiron/sqlx"
)

type Bot interface {
	GetSubscriberByChatId(ChatId int64) (core.Subscriber, error)
	CreateSubscriber(ChatId int64) error
	ActivateSubscriber(chatId int64) error
}

type Repository struct {
	Bot
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Bot: NewBotPostgres(db),
	}
}
