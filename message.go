package qbot

type Message struct {
	Date       string
	FromUserId int64
	MessageId  int
	ChatId     int64
	Text       string
	Json       string
	Mailing    int
	IsUnknown  string
}
