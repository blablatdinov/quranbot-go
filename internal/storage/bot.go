package storage

import (
	"github.com/blablatdinov/quranbot-go/internal/core"
	"github.com/jmoiron/sqlx"
)

type BotPostgres struct {
	db *sqlx.DB
}

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

func (r *BotPostgres) CreateSubscriber(chatId int64) error {
	query := `INSERT INTO bot_init_subscriber (tg_chat_id, is_active, day) VALUES
	($1, 't', 2)`
	_, err := r.db.Exec(query, chatId)
	return err
}

func (r *BotPostgres) ActivateSubscriber(chatId int64) error {
	query := `UPDATE bot_init_subscriber SET is_active = 't' WHERE tg_chat_id = $1`
	_, err := r.db.Exec(query, chatId)
	return err
}
