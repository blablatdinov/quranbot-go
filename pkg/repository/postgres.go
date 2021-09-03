package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewPostgres() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=almazilaletdinov dbname=qbot sslmode=disable"))
	if err != nil {
		return &sqlx.DB{}, err
	}
	return db, nil
}
