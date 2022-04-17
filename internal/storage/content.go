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
	var ayats []core.Ayat
	query := "SELECT id, ayat FROM content_ayats a" +
		"INNER JOIN content_sura s on s.id = a.sura_id" +
		"WHERE s.number = $1"
	if err := r.db.Select(ayats, query, suraNum); err != nil {
		return []core.Ayat{}, err
	}
	return ayats, nil
}
