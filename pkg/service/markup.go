package service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func GetDefaultKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Подкасты"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Время намаза"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Избранное"),
			tgbotapi.NewKeyboardButton("Найти аят"),
		),
	)
}
