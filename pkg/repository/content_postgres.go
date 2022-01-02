package repository

import (
	"errors"
	"fmt"
	"log"
	"qbot"

	"github.com/jmoiron/sqlx"
)

type ContentPostgres struct {
	db *sqlx.DB
}

func NewContentPostgres(db *sqlx.DB) *ContentPostgres {
	return &ContentPostgres{db: db}
}

const (
	defaultAyatFields = "a.id, a.ayat, a.content, s.number as sura_number, s.link as sura_link"
)

func (r *ContentPostgres) GetAyatByMailingDay(mailingDay int) (qbot.Ayat, error) {
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

func (r *ContentPostgres) GetAyatsBySuraNum(suraNum int) ([]qbot.Ayat, error) {
	var ayats []qbot.Ayat
	query := `select 
		a.id,
		a.ayat,
		a.content,
		s.number as sura_number,
		s.link as sura_link
	from content_ayat a 
    inner join content_sura s on a.sura_id = s.id 
	where s.number = $1`
	if err := r.db.Select(&ayats, query, suraNum); err != nil {
		return []qbot.Ayat{}, nil
	}
	return ayats, nil
}

func (r *ContentPostgres) GetFavoriteAyats(chatId int64) ([]qbot.Ayat, error) {
	var ayats []qbot.Ayat
	query := `select 
			a.id,
			a.ayat,
			a.content,
			s.number as sura_number,
			s.link as sura_link
		from bot_init_subscriber sub
		inner join bot_init_subscriber_favourite_ayats fa on sub.id = fa.subscriber_id
		inner join content_ayat a on fa.ayat_id = a.id
		inner join content_sura s on a.sura_id = s.id
		where sub.tg_chat_id = $1
		order by a.id`
	if err := r.db.Select(&ayats, query, chatId); err != nil {
		return []qbot.Ayat{}, err
	}
	return ayats, nil
}

func (r *ContentPostgres) GetRandomPodcast() (qbot.Podcast, error) {
	var podcast qbot.Podcast
	query := `
		select
			cf.tg_file_id as tg_file_id,
			cf.link_to_file as link_to_file
		from content_podcast cp
		inner join content_file cf on cf.id = cp.audio_id
		order by random()
		limit 1`
	if err := r.db.Get(&podcast, query); err != nil {
		return qbot.Podcast{}, err
	}
	return podcast, nil
}

func (r *ContentPostgres) AyatIsFavorite(chatId int64, ayatId int) bool {
	// TODO: rename bot_init_subscriber_favourite_ayats table
	var count int
	query := `
	select
		count(*)
	from bot_init_subscriber_favourite_ayats fa
	inner join content_ayat a on a.id = fa.ayat_id
	inner join bot_init_subscriber sub on sub.id = fa.subscriber_id
	where sub.tg_chat_id = $1 and a.id = $2`
	if err := r.db.Get(&count, query, chatId, ayatId); err != nil {
		return false
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

func (r *ContentPostgres) RemoveFromFavorite(chatId int64, ayatId int) error {
	var favoriteAyatIdForRemove int
	query := `
	select 
	    fa.id
	from content_ayat a
	inner join bot_init_subscriber_favourite_ayats fa on fa.ayat_id = a.id
	inner join bot_init_subscriber sub on sub.id = fa.subscriber_id
	where sub.tg_chat_id = $1 and fa.ayat_id = $2`
	err := r.db.Get(&favoriteAyatIdForRemove, query, chatId, ayatId)
	if err != nil {
		return err
	}
	query = `
	delete
	from bot_init_subscriber_favourite_ayats
	where id = $1`
	_, err = r.db.Exec(query, favoriteAyatIdForRemove)
	if err != nil {
		return err
	}
	log.Printf("fa %d deleted\n", favoriteAyatIdForRemove)
	return nil
}

func (r *ContentPostgres) GetAyatById(chatId int64, ayatId int) (qbot.Ayat, error) {
	var ayat qbot.Ayat
	query := `
	select 
		 a.id,
		 a.ayat,
		 a.content,
		 s.number as sura_number,
		 s.link as sura_link
	from content_ayat a
	inner join content_sura s on a.sura_id = s.id
	where a.id = $1`
	if err := r.db.Get(&ayat, query, ayatId); err != nil {
		return qbot.Ayat{}, err
	}
	ayat.IsFavorite = r.AyatIsFavorite(chatId, ayatId)
	return ayat, nil
}

func (r *ContentPostgres) AddToFavorite(chatId int64, ayatId int) error {
	var subscriberId int64
	if r.AyatIsFavorite(chatId, ayatId) {
		return errors.New("ayat is already favorite")
	}
	query := fmt.Sprintf("select id from bot_init_subscriber where tg_chat_id=%d", chatId)
	if err := r.db.Get(&subscriberId, "select id from bot_init_subscriber where tg_chat_id=$1", chatId); err != nil {
		return err
	}
	query = `
	insert into bot_init_subscriber_favourite_ayats
	(subscriber_id, ayat_id)
	values
	($1, $2)
	returning id`
	_, err := r.db.Exec(query, subscriberId, ayatId)
	return err
}

func (r *ContentPostgres) GetAdjacentAyats(chatId int64, ayatId int) ([]qbot.Ayat, error) {
	var ayats []qbot.Ayat
	query := fmt.Sprintf(`
	select %s
	from (
		select 
			%s,
			lag(a.id) over (order by a.id asc) as prev,
			lead(a.id) over (order by a.id asc) as next
		from content_ayat a
			inner join bot_init_subscriber_favourite_ayats fa on a.id = fa.ayat_id
			inner join bot_init_subscriber sub on sub.id = fa.subscriber_id
			inner join content_sura s on s.id = a.sura_id
		where sub.tg_chat_id = $1
		) x
	where $2 IN (id, prev, next)`, defaultAyatFields, defaultAyatFields)
	err := r.db.Select(&ayats, query, chatId, ayatId)
	return ayats, err
}

// GetMorningContentForTodayMailing достать из базы данных контент для утренней рассылки
func (r *ContentPostgres) GetMorningContentForTodayMailing() ([]qbot.MailingContent, error) {
	var contentForMailing []qbot.MailingContent
	query := `
	select
		s.tg_chat_id,
		STRING_AGG(
			'__' || sura.number::character varying || ':' || a.ayat || ')__ ' || a .content || '\n\n',
			''
			order by a.id
		) as content,
		STRING_AGG(sura.link, '|' order by a.id) as link
	from bot_init_subscriber as s
	inner join content_morningcontent as mc on s.day=mc.day
	inner join content_ayat as a on a.one_day_content_id=mc.id
	inner join content_sura as sura on a.sura_id=sura.id
	where s.is_active = 'true'
	group by s.tg_chat_id`
	err := r.db.Select(&contentForMailing, query)
	return contentForMailing, err
}

func (r *ContentPostgres) UpdateDaysForSubscribers(chatIds []int64) error {
	query := fmt.Sprintf(`
	update bot_init_subscriber
	set day = day + 1
	%s`, GenerateConditionForUpdatingSubscribers(chatIds))
	_, err := r.db.Exec(query)
	return err
}
