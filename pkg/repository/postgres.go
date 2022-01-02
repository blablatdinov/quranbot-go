package repository

import (
	"github.com/jmoiron/sqlx"
)

func NewPostgres(databaseUrl string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseUrl)
	if err != nil {
		return &sqlx.DB{}, err
	}
	return db, nil
}
