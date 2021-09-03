package main

import (
	"log"
	"qbot/pkg/repository"
	"qbot/pkg/service"
	"qbot/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func main() {
	botApi, err := tgbotapi.NewBotAPI("452230948:AAFBRxigSQIg1PZJTjrY6c3OHFvTnKDP_AA")
	if err != nil {
		log.Panic(err)
	}
	botApi.Debug = false

	db, err := repository.NewPostgres()
	if err != nil {
		log.Println(err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	bot := telegram.NewBot(botApi, services)
	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}
