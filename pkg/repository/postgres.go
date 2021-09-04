package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewPostgres(databaseUrl string) (*sqlx.DB, error) {
	//db, err := sqlx.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=almazilaletdinov dbname=qbot sslmode=disable"))
	fmt.Println(databaseUrl)
	db, err := sqlx.Open("postgres", databaseUrl)
	if err != nil {
		return &sqlx.DB{}, err
	}
	return db, nil
}
