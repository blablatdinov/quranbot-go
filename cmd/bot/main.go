package main

import (
	"log"

	"github.com/blablatdinov/quranbot-go/internal/service"

	"github.com/blablatdinov/quranbot-go/internal/storage"
	"github.com/blablatdinov/quranbot-go/pkg/telegramsdk"
	_ "github.com/lib/pq"
)

func main() {
	bot := telegramsdk.NewBot("452230948:AAFvAXqcuK8xhw1gfGnxlp6zzWQaR9qK7hw")
	updatesChan := bot.GetUpdatesChan()
	db, err := storage.NewPostgres("postgres://almazilaletdinov@localhost:5432/qbot?sslmode=disable")
	repos := storage.NewRepository(db)
	services := service.NewService(repos)
	if err != nil {
		log.Fatal(err.Error())
	}
	for message := range updatesChan {
		if message.Text == "/start" {
			text, err := services.GetOrCreateSubscriber(message.Chat.Id)
			if err != nil {
				log.Fatal(err.Error())
			}
			bot.SendMessage(message.Chat.Id, text)
		}
	}
}
