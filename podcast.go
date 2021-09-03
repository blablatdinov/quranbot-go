package qbot

import "database/sql"

type Podcast struct {
	TgFileId   sql.NullString `db:"tg_file_id"`
	LinkToFile string         `db:"link_to_file"`
}
