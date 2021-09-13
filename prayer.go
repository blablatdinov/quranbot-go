package qbot

import "time"

type Prayer struct {
	Id   int       `db:"id"`
	Time time.Time `db:"time"`
	Date time.Time `db:"date"`
	Name string    `db:"name"`

	CityName string `db:"city_name"`
}
