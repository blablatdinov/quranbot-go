package storage

import (
	"strconv"

	"github.com/blablatdinov/quranbot-go/internal/core"
	"github.com/jmoiron/sqlx"
)

type BotPostgres struct {
	db *sqlx.DB
}

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

func NewBotPostgres(db *sqlx.DB) *BotPostgres {
	return &BotPostgres{
		db: db,
	}
}

func (r *BotPostgres) GetSubscriberByChatId(chatId int64) (core.Subscriber, error) {
	var subscriber core.Subscriber
	query := `SELECT
		id,
		tg_chat_id,
		is_active,
		day,
		step
	FROM bot_init_subscriber
	WHERE tg_chat_id = $1`
	if err := r.db.Get(&subscriber, query, chatId); err != nil {
		return subscriber, err
	}
	return subscriber, nil
}

func (r *BotPostgres) CreateSubscriber(chatId int64, referralCode string) error {
	refererId, err := strconv.Atoi(referralCode)
	if err != nil && referralCode == "0" {
		query := `INSERT INTO bot_init_subscriber (tg_chat_id, is_active, day) VALUES
		($1, 't', 2)`
		_, err := r.db.Exec(query, chatId)
		return err
	} else {
		query := `INSERT INTO bot_init_subscriber (tg_chat_id, is_active, day, referer_id) VALUES
		($1, 't', 2, $2)`
		_, err := r.db.Exec(query, chatId, refererId)
		return err
	}
}

func (r *BotPostgres) ActivateSubscriber(chatId int64) error {
	query := `UPDATE bot_init_subscriber SET is_active = 't' WHERE tg_chat_id = $1`
	_, err := r.db.Exec(query, chatId)
	return err
}
