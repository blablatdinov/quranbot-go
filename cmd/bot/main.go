package main

import (
	"github.com/blablatdinov/quranbot-go/pkg/telegramsdk"
)

func main() {
	bot := telegramsdk.NewBot("452230948:AAFvAXqcuK8xhw1gfGnxlp6zzWQaR9qK7hw")
	updatesChan := bot.GetUpdatesChan()
	for message := range updatesChan {
		bot.SendMessage(message.Chat.Id, message.Text)
	}
}
