package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"qbot/pkg/repository"
	"qbot/pkg/service"
	"qbot/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	botApi, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	debugMode := os.Getenv("DEBUG") == "true"
	botApi.Debug = debugMode
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalln("Set DATABASE_URL enviroment variable")
	}
	db, err := repository.NewPostgres(databaseUrl)
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
