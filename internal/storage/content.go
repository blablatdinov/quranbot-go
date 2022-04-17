package storage

import (
	"fmt"
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
	fmt.Println(suraNum)
	query := "SELECT a.id as id, a.ayat as ayat FROM content_ayat a " +
		"INNER JOIN content_sura s on s.id = a.sura_id " +
		"WHERE s.number = $1"
	if err := r.db.Select(&ayats, query, suraNum); err != nil {
		return []core.Ayat{}, err
	}
	return ayats, nil
}
