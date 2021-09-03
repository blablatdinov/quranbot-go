package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"qbot"
)

type PrayerPostgres struct {
	db *sqlx.DB
}

func NewPrayerPostgres(db *sqlx.DB) *PrayerPostgres {
	return &PrayerPostgres{db}
}

func (r *PrayerPostgres) SubscriberHasCity(chatId int64) (bool, error) {
	var cityId sql.NullInt16
	var x int16 = 1
	query := "select city_id from bot_init_subscriber where tg_chat_id = $1"
	err := r.db.Get(&cityId, query, chatId)
	if err != nil {
		return false, err
	}
	if cityId.Int16 > x {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *PrayerPostgres) GetPrayer(chatId int64) ([]qbot.Prayer, error) {
	var ayats []qbot.Prayer
	query := `
		select
			city.name as city_name,
			day.date,
			p.time
		from prayer_prayer p
		inner join prayer_city city on city.id = p.city_id
		inner join bot_init_subscriber sub on city.id = sub.city_id
		inner join prayer_day day on p.day_id = day.id
		where sub.tg_chat_id = $1`
	err := r.db.Select(&ayats, query, chatId)
	return ayats, err
}
