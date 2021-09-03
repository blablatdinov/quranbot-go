package qbot

import "fmt"

type Ayat struct {
	Id int `db:"id"`
	//AdditionalContent string `db:"additional_content"`
	Content    string `db:"content"`
	Arab_text  string `db:"arab_text"`
	Trans      string `db:"trans"`
	IsFavorite bool   `db:"is_favorite"`

	Sura       int    `db:"sura_id"`
	SuraLink   string `db:"sura_link"`
	SuraNumber int    `db:"sura_number"`

	Ayat            string `db:"ayat"`
	Html            string `db:"html"`
	Audio           int    `db:"audio_id"`
	One_day_content int    `db:"one_day_content_id"`
}

func (a *Ayat) GetSuraAyatNum() string {
	return fmt.Sprintf("%d:%s", a.SuraNumber, a.Ayat)
}
