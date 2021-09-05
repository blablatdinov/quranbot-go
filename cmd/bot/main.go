package main

import (
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"qbot/pkg/repository"
	"qbot/pkg/service"
	"qbot/pkg/telegram"
	"time"
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
	timezone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalln("Set DATABASE_URL enviroment variable")
	}
	goCron := gocron.NewScheduler(timezone)
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	bot := telegram.NewBot(botApi, services, goCron)
	if err := bot.StartJobs(); err != nil {
		log.Fatal(err)
	}
	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}
