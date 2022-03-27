package main

import (
	"fmt"
	"log"

	"github.com/blablatdinov/quranbot-go/internal/service"

	"github.com/blablatdinov/quranbot-go/internal/storage"
	"github.com/blablatdinov/quranbot-go/pkg/telegramsdk"
	_ "github.com/lib/pq"
)

func main() {
	bot := telegramsdk.NewBot("452230948:AAFvAXqcuK8xhw1gfGnxlp6zzWQaR9qK7hw")
	// updatesChan := bot.GetUpdatesChan()
	updatesChan, err := bot.RunWebhookServer("https://3a1d-87-117-185-236.ngrok.io")
	if err != nil {
		log.Fatal(err.Error())
	}
	db, err := storage.NewPostgres("postgres://almazilaletdinov@localhost:5432/qbot?sslmode=disable")
	repos := storage.NewRepository(db)
	services := service.NewService(repos)
	if err != nil {
		log.Fatal(err.Error())
	}
	for message := range updatesChan {
		fmt.Printf("message text: %s\n", message.Text)
		if len(message.Text) > 5 && message.Text[:6] == "/start" {
			referralCode := "0"
			if len(message.Text) > 6 {
				referralCode = message.Text[7:]
			}
			text, err := services.GetOrCreateSubscriber(message.Chat.Id, referralCode)
			if err != nil {
				log.Fatal(err.Error())
			}
			bot.SendMessage(message.Chat.Id, text)
		}
	}
}
