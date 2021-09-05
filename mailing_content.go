package qbot

type MailingContent struct {
	ChatId  int64  `db:"tg_chat_id"`
	Content string `db:"content"`
	Link    string `db:"link"`
}
