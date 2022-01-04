package main

import (
	"flag"
	"log"
	"os"
	"qbot/pkg/repository"
	"qbot/pkg/service"
	"qbot/pkg/telegram"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

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
	databaseUrl := os.Getenv("GO_DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalln("Set GO_DATABASE_URL enviroment variable")
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
	adminsList := []int64{}
	for _, adminIdStr := range strings.Split(os.Getenv("ADMINS"), ",") {
		adminId, err := strconv.Atoi(adminIdStr)
		if err != nil {
			log.Fatal("Check ADMINS env variable")
		}
		adminsList = append(adminsList, int64(adminId))
	}
	bot := telegram.NewBot(botApi, services, gocron.NewScheduler(time.UTC), adminsList)
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
	} else if flag.Args()[0] == "cron" {
		var wg sync.WaitGroup
		wg.Add(1)
		if err := bot.StartJobs(); err != nil {
			log.Fatal(err)
		}
		wg.Wait()
	} else {
		log.Fatal("Command not found")
	}
}
