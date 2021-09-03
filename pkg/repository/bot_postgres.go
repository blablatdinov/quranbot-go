package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"qbot"
)

type BotPostgres struct {
	db *sqlx.DB
}

func NewBotPostgres(db *sqlx.DB) *BotPostgres {
	return &BotPostgres{
		db: db,
	}
}

func (r *BotPostgres) GetSubscriberByChatId(chatId int64) (qbot.Subscriber, error) {
	var subscriber qbot.Subscriber
	query := fmt.Sprintf(
		`SELECT 
			id,
			tg_chat_id,
			is_active,
			day
		FROM bot_init_subscriber AS s 
		WHERE s.tg_chat_id=%d`,
		chatId,
	)
	if err := r.db.Get(&subscriber, query); err != nil {
		return qbot.Subscriber{}, err
	}
	return subscriber, nil
}

func (r *BotPostgres) CreateSubscriber(chatId int64) (qbot.Subscriber, error) {
	var subscriberChatId int64
	query := "INSERT INTO bot_init_subscriber (tg_chat_id, is_active, day) VALUES ($1, $2, $3) RETURNING tg_chat_id"
	row := r.db.QueryRow(query, chatId, true, 2)
	err := row.Scan(&subscriberChatId)
	if err != nil {
		return qbot.Subscriber{}, err
	}
	subscriber, err := r.GetSubscriberByChatId(subscriberChatId)
	if err != nil {
		return qbot.Subscriber{}, err
	}
	return subscriber, nil
}

func (r *BotPostgres) GetOrCreateSubscriber(chatId int64) (qbot.Subscriber, bool, error) {
	subscriber, err := r.GetSubscriberByChatId(chatId)
	if err == nil {
		return subscriber, false, nil
	} else if err.Error() != "sql: no rows in result set" {
		return qbot.Subscriber{}, false, err
	}

	subscriber, err = r.CreateSubscriber(chatId)
	if err == nil {
		return subscriber, true, nil
	}

	return qbot.Subscriber{}, false, err
}

func (r *BotPostgres) SetSubscriberState(chatId int64, step string) error {
	query := "update bot_init_subscriber set step = $2 where tg_chat_id = $1"
	_, err := r.db.Exec(query, chatId, step)
	return err
}

func (r *BotPostgres) GetSubscriberState(chatId int64) (string, error) {
	var state string
	query := "select step from bot_init_subscriber where tg_chat_id = $1"
	err := r.db.Get(&state, query, chatId)
	return state, err
}
