package storage

import (
	"github.com/blablatdinov/quranbot-go/internal/core"
	"github.com/jmoiron/sqlx"
)

type ContentRepository struct {
	db *sqlx.DB
}

func NewContentPostgres(db *sqlx.DB) *ContentRepository {
	return &ContentRepository{db}
}

func (r *ContentRepository) GetAyatsBySuraNum(suraNum int) ([]core.Ayat, error) {
	return []core.Ayat{}, nil
}
