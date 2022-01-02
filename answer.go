package qbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// Answer структура для ответа пользователю
type Answer struct {
	ChatId   int64  `db:"tg_chat_id"`
	Content  string `db:"content"`
	Keyboard tgbotapi.InlineKeyboardMarkup
}

func (a *Answer) HasKeyboard() bool {
	if len(a.Keyboard.InlineKeyboard) == 0 {
		return false
	}
	return true
}
