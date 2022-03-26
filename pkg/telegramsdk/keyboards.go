package telegramsdk

import (
	"encoding/json"
)

type ReplyKeyboardMarkup struct {
	Keyboard [][]ReplyKeyboardButton `json:"keyboard"`
}

type ReplyKeyboardButton struct {
	Text string `json:"text"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

func NewInlineKeyboardMarkup(buttons [][]InlineKeyboardButton) (string, error) {
	defaultKeyboard := InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
	keyboardJson, err := json.Marshal(defaultKeyboard)
	if err != nil {
		return "", err
	}
	return string(keyboardJson), err
}
