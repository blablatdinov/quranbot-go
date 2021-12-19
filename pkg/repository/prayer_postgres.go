package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"qbot"
	"time"
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

// GetPrayer получить время намаза для пользователя по идентификатору чата и дате
func (r *PrayerPostgres) GetPrayer(chatId int64, date time.Time) ([]qbot.Prayer, error) {
	var prayers []qbot.Prayer
	query := `
		select
			p.id,
			city.name as city_name,
			day.date,
			p.time
		from prayer_prayer p
		inner join prayer_city city on city.id = p.city_id
		inner join bot_init_subscriber sub on city.id = sub.city_id
		inner join prayer_day day on p.day_id = day.id
		where sub.tg_chat_id = $1 and day.date = $2`
	err := r.db.Select(&prayers, query, chatId, date.Format("01-02-2006"))
	return prayers, err
}

func (r *PrayerPostgres) CreatePrayerAtUserGroup() (int, error) {
	var prayerAtUserGroupId int
	query := `insert into prayer_prayeratusergroup DEFAULT VALUES returning id`
	row := r.db.QueryRow(query)
	if err := row.Scan(&prayerAtUserGroupId); err != nil {
		return 0, err
	}
	return prayerAtUserGroupId, nil
}

func (r *PrayerPostgres) GeneratePrayerForUser(chatId int64, prayers []qbot.Prayer) ([]qbot.PrayerAtUser, error) {
	prayerAtUserGroupId, err := r.CreatePrayerAtUserGroup()
	if err != nil {
		return []qbot.PrayerAtUser{}, err
	}
	var subscriberId int
	query := "select id from bot_init_subscriber where tg_chat_id = $1"
	if err := r.db.Get(&subscriberId, query, chatId); err != nil {
		return []qbot.PrayerAtUser{}, err
	}

	query = `insert into prayer_prayeratuser
	(is_read, prayer_id, prayer_group_id, subscriber_id)
	values` +
		fmt.Sprintf("('f', %d, %d, %d),", prayers[0].Id, prayerAtUserGroupId, subscriberId) +
		fmt.Sprintf("('f', %d, %d, %d),", prayers[2].Id, prayerAtUserGroupId, subscriberId) +
		fmt.Sprintf("('f', %d, %d, %d),", prayers[3].Id, prayerAtUserGroupId, subscriberId) +
		fmt.Sprintf("('f', %d, %d, %d),", prayers[4].Id, prayerAtUserGroupId, subscriberId) +
		fmt.Sprintf("('f', %d, %d, %d)", prayers[5].Id, prayerAtUserGroupId, subscriberId)
	_, err = r.db.Exec(query)
	if err != nil {
		return []qbot.PrayerAtUser{}, err
	}
	var prayersAtUser []qbot.PrayerAtUser
	query = "select id, is_read from prayer_prayeratuser where prayer_group_id = $1"
	if err = r.db.Select(&prayersAtUser, query, prayerAtUserGroupId); err != nil {
		return []qbot.PrayerAtUser{}, err
	}
	return prayersAtUser, nil
}

func (r *PrayerPostgres) GetPrayerForUser(chatId int64, prayers []qbot.Prayer) ([]qbot.PrayerAtUser, error) {
	var prayersAtUser []qbot.PrayerAtUser
	query := `
	select 
		pp.id, 
		is_read 
	from prayer_prayeratuser as pp
	inner join bot_init_subscriber bis on bis.id = pp.subscriber_id
	where bis.tg_chat_id = $1 and (pp.prayer_id = $2 or pp.prayer_id = $3 or pp.prayer_id = $4 or pp.prayer_id = $5 or pp.prayer_id = $6)`
	if err := r.db.Select(&prayersAtUser, query, chatId, prayers[0].Id, prayers[2].Id, prayers[3].Id, prayers[4].Id, prayers[5].Id); err != nil {
		return []qbot.PrayerAtUser{}, err
	}
	return prayersAtUser, nil
}

func (r *PrayerPostgres) GetOrCreatePrayerForUser(chatId int64, prayers []qbot.Prayer) ([]qbot.PrayerAtUser, error) {
	prayersAtUser, err := r.GetPrayerForUser(chatId, prayers)
	if err != nil {
		return []qbot.PrayerAtUser{}, err
	}
	if len(prayersAtUser) == 0 {
		prayersAtUser, err = r.GeneratePrayerForUser(chatId, prayers)
		if err != nil {
			return []qbot.PrayerAtUser{}, err
		}
		return prayersAtUser, nil
	}
	return prayersAtUser, nil
}

func (r *PrayerPostgres) ChangePrayerStatus(prayerAtUserId int, status bool) error {
	query := "update prayer_prayeratuser set is_read = $1 where id = $2"
	_, err := r.db.Exec(query, status, prayerAtUserId)
	return err
}

func (r *PrayerPostgres) GetPrayersAtUserByGroupId(prayersAtUserGroupId int) ([]qbot.PrayerAtUser, error) {
	var prayers []qbot.PrayerAtUser
	query := "select id, is_read from prayer_prayeratuser where prayer_group_id = $1"
	if err := r.db.Select(&prayers, query, prayersAtUserGroupId); err != nil {
		return []qbot.PrayerAtUser{}, err
	}
	return prayers, nil
}

func (r *PrayerPostgres) GetPrayersAtUserByOnePrayerId(prayersAtUserId int) ([]qbot.PrayerAtUser, error) {
	var prayers []qbot.PrayerAtUser
	query := `
	select
		p.id,
		is_read
	from prayer_prayeratuser as p
         inner join prayer_prayeratusergroup pp on p.prayer_group_id = pp.id
	where pp.id = (select prayer_group_id from prayer_prayeratuser where id=$1)`
	if err := r.db.Select(&prayers, query, prayersAtUserId); err != nil {
		return []qbot.PrayerAtUser{}, err
	}
	return prayers, nil
}

func (r *PrayerPostgres) GetCityByName(cityName string) (qbot.City, error) {
	var city qbot.City
	query := "select id, name from prayer_city where name = $1"
	err := r.db.Get(&city, query, cityName)
	return city, err
}

func (r *PrayerPostgres) ChangeCity(chatId int64, cityId int) error {
	query := "update bot_init_subscriber set city_id = $1 where tg_chat_id = $2"
	_, err := r.db.Exec(query, cityId, chatId)
	return err
}

// GetSubscriberWithCityChatIds получить массив идентификаторов пользователей, у которых установлен город
func (r *PrayerPostgres) GetSubscriberWithCityChatIds() ([]int64, error) {
	var subscriberWithCityChatIds []int64
	query := `
		select
			b.tg_chat_id
		from bot_init_subscriber b
		where b.city_id is not null and is_active = true`
	err := r.db.Select(&subscriberWithCityChatIds, query)
	if err != nil {
		return []int64{}, err
	}
	return subscriberWithCityChatIds, nil
}
