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
}
