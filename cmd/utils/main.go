package main

import (
	"flag"
	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
	"qbot/pkg/repository"
	"qbot/pkg/service"
	"qbot/pkg/telegram"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func initialize() (*tgbotapi.BotAPI, *sqlx.DB) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	botApi, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	botApi.Debug = false
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalln("Set DATABASE_URL enviroment variable")
	}
	db, err := repository.NewPostgres(databaseUrl)
	if err != nil {
		log.Println(err.Error())
	}
	return botApi, db
}

func main() {
	botApi, db := initialize()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	bot := telegram.NewBot(botApi, services, gocron.NewScheduler(time.UTC))
	flag.Parse()
	if flag.Args()[0] == "check_subscribers" {
		if err := bot.CheckSubscribers(); err != nil {
			log.Fatal(err)
		}
	} else if flag.Args()[0] == "send_content" {
		if err := bot.SendMorningContent(); err != nil {
			log.Fatal(err)
		}
	} else if flag.Args()[0] == "send_prayers" {
		if err := bot.SendPrayerTimes(); err != nil {
			log.Fatal(err)
		}
	}
}
