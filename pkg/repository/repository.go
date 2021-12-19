package repository

import (
	"github.com/jmoiron/sqlx"
	"qbot"
	"time"
)

type Bot interface {
	GetOrCreateSubscriber(chatId int64) (qbot.Subscriber, bool, error)
	SetSubscriberState(chatId int64, step string) error
	GetSubscriberState(chatId int64) (string, error)
	GetAyatByMailingDay(mailingDay int) (qbot.Ayat, error)
	GetActiveSubscribers() ([]qbot.Subscriber, error)
	DeactivateSubscribers([]int64) error
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
	GetMorningContentForTodayMailing() ([]qbot.MailingContent, error)
	UpdateDaysForSubscribers([]int64) error
}

type Prayer interface {
	GetPrayer(chatId int64, date time.Time) ([]qbot.Prayer, error)
	SubscriberHasCity(chatId int64) (bool, error)
	GeneratePrayerForUser(chatId int64, prayers []qbot.Prayer) ([]qbot.PrayerAtUser, error)
	GetOrCreatePrayerForUser(chatId int64, prayers []qbot.Prayer) ([]qbot.PrayerAtUser, error)
	ChangePrayerStatus(prayerAtUserId int, status bool) error
	GetPrayersAtUserByGroupId(prayersAtUserGroupId int) ([]qbot.PrayerAtUser, error)
	GetPrayersAtUserByOnePrayerId(prayersAtUserId int) ([]qbot.PrayerAtUser, error)
	GetCityByName(cityName string) (qbot.City, error)
	ChangeCity(chatId int64, cityId int) error
	GetSubscriberWithCityChatIds() ([]int64, error)
}

type Repository struct {
	Bot
	Content
	Prayer
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Bot:     NewBotPostgres(db),
		Content: NewContentPostgres(db),
		Prayer:  NewPrayerPostgres(db),
	}
}
