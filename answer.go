package qbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// Answer структура для ответа пользователю
type Answer struct {
	ChatId   int64
	Content  string
	Keyboard tgbotapi.InlineKeyboardMarkup
}
