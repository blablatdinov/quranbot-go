package main

import (
	"fmt"
	"log"

	"github.com/blablatdinov/quranbot-go/pkg/telegramsdk"
)

func main() {
	bot := telegramsdk.NewBot("452230948:AAFvAXqcuK8xhw1gfGnxlp6zzWQaR9qK7hw")
	message, err := bot.SendMessage(358610865, "bot")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(message.Result.Date)
	keyboard, _ := telegramsdk.NewInlineKeyboardMarkup([][]telegramsdk.InlineKeyboardButton{
		{
			{Text: "1", CallbackData: "1"},
			{Text: "2", CallbackData: "2"},
		},
		{
			{Text: "3", CallbackData: "3"},
			{Text: "4", CallbackData: "4"},
		},
	})
	message, err = bot.SendMessageWithKeyboard(358610865, "message with inline keyboard", keyboard)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(message.Result.Date)
}
