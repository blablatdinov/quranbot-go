package core

import "database/sql"

type Subscriber struct {
	Id       int            `db:"id"`
	ChatId   int64          `db:"tg_chat_id"`
	IsActive bool           `db:"is_active"`
	Day      int            `db:"day"`
	Step     sql.NullString `db:"step"`
}
