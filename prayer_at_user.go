package qbot

type PrayerAtUser struct {
	Id     int  `db:"id"`
	IsRead bool `db:"is_read"`
}
