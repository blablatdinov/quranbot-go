package repository

import (
	"fmt"
	"qbot"
	"strings"
	"time"

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

func (r *BotPostgres) GetAyatByMailingDay(mailingDay int) (qbot.Ayat, error) {
	var ayat qbot.Ayat
	query := `
		select 
			a.id,
		    a.content,
		    a.arab_text,
		    a.trans,
		    a.sura_id,
		    s.link as sura_link,
		    a.ayat,
		    a.html,
		    a.audio_id,
		    a.one_day_content_id
		from content_ayat a
		inner join content_morningcontent cm on a.one_day_content_id = cm.id
		inner join content_sura s on a.sura_id = s.id
		where cm.day = $1`
	if err := r.db.Get(&ayat, query, mailingDay); err != nil {
		return qbot.Ayat{}, err
	}
	return ayat, nil
}

func (r *BotPostgres) GetActiveSubscribers() ([]qbot.Subscriber, error) {
	var subscribers []qbot.Subscriber
	query := `
	select 
		tg_chat_id
	from bot_init_subscriber
	where is_active = 't'
	`
	err := r.db.Select(&subscribers, query)
	return subscribers, err
}

// GetActiveSubscribersCount получить кол-во активных подписчиков
func (r *BotPostgres) GetSubscribersCount(param string) (int, error) {
	var result int
	var query string
	if param == "total" {
		query = "select count(*) from bot_init_subscriber"
	} else if param == "active" {
		query = "select count(*) from bot_init_subscriber where is_active='t'"
	} else {
		return 0, fmt.Errorf("unsupported param: %s", param)
	}

	row := r.db.QueryRow(query)
	if err := row.Scan(&result); err != nil {
		return 0, err
	}
	return result, nil
}

func GenerateConditionForUpdatingSubscribers(chatIds []int64) string {
	result := "WHERE "
	var or string
	for i, chatId := range chatIds {
		if i == len(chatIds)-1 {
			or = ""
		} else {
			or = " OR "
		}
		result += fmt.Sprintf("tg_chat_id=%d%s", chatId, or)
	}
	return result
}

func (r *BotPostgres) DeactivateSubscribers(chatIds []int64) error {
	query := fmt.Sprintf(`
		UPDATE bot_init_subscriber
		SET is_active = 'f'
	%s`, GenerateConditionForUpdatingSubscribers(chatIds))
	_, err := r.db.Exec(query)
	return err
}

func (r *BotPostgres) CreateSubscriberActions(chatIds []int64, action string) error {
	subscriberIds, err := r.getSubsriberIdsFromChatIds(chatIds)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO bot_init_subscriberaction
		(action, subscriber_id, date_time)
		VALUES
	`
	dateTime := time.Now().Format("2006-02-01 15:04:05.999999999") + " +03:00"
	for _, subscriberId := range subscriberIds {
		query += fmt.Sprintf("('%s', %d, '%s')", action, subscriberId, dateTime)
	}
	_, err = r.db.Exec(query)
	return err
}

func (r *BotPostgres) getSubsriberIdsFromChatIds(chatIds []int64) ([]int, error) {
	var result []int
	condition := GenerateConditionForUpdatingSubscribers(chatIds)
	query := fmt.Sprintf("SELECT id from bot_init_subscriber %s", condition)
	if err := r.db.Select(&result, query); err != nil {
		return []int{}, err
	}
	return result, nil
}

func (r *BotPostgres) SaveMessage(message qbot.Message) error {
	if message.Mailing == 0 {
		query := `
			insert into bot_init_message
			(date, from_user_id, message_id, chat_id, text, json, is_unknown)
			values
			($1, $2, $3, $4, $5, $6, $7)
		`
		_, err := r.db.Exec(query, message.Date, message.FromUserId, message.MessageId, message.ChatId, message.Text, message.Json, message.IsUnknown)
		return err
	} else {
		query := `
			insert into bot_init_message
			(date, from_user_id, message_id, chat_id, text, json, is_unknown, mailing_id)
			values
			($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err := r.db.Exec(query, message.Date, message.FromUserId, message.MessageId, message.ChatId, message.Text, message.Json, message.IsUnknown, message.Mailing)
		return err
	}
}

func (r *BotPostgres) BulkSaveMessages(messages []qbot.Message) error {
	query := `
		insert into bot_init_message
		(date, from_user_id, message_id, chat_id, text, json, mailing_id, is_unknown)
		values
	`
	valuesArray := []string{}
	for _, message := range messages {
		valuesArray = append(valuesArray, fmt.Sprintf(
			"('%s'::timestamptz, %d, %d, %d, '%s', '%s', %d, '%s')",
			message.Date,
			message.FromUserId,
			message.MessageId,
			message.ChatId,
			message.Text,
			message.Json,
			message.Mailing,
			message.IsUnknown,
		))
	}
	values := strings.Join(valuesArray, ",")
	_, err := r.db.Exec(query + values)
	if err != nil {
		return err
	}
	return nil
}

func (r *BotPostgres) CreateMailing() (int, error) {
	var mailingId int
	query := "insert into bot_init_mailing default values returning id"
	row := r.db.QueryRow(query)
	err := row.Scan(&mailingId)
	if err != nil {
		return 0, err
	}
	return mailingId, nil
}
