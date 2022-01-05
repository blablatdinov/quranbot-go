package main

import (
	"log"
	"os"
	"qbot/pkg/repository"
	"qbot/pkg/service"
	"qbot/pkg/telegram"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка dotenv файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	// Инициализация бота
	botApi, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	debugMode := os.Getenv("DEBUG") == "true"
	botApi.Debug = debugMode
	databaseUrl := os.Getenv("GO_DATABASE_URL")
	// Инициализация списка администраторов
	adminsList := []int64{}
	for _, adminIdStr := range strings.Split(os.Getenv("ADMINS"), ",") {
		adminId, err := strconv.Atoi(adminIdStr)
		if err != nil {
			log.Fatal("Check ADMINS env variable")
		}
		adminsList = append(adminsList, int64(adminId))
	}
	if databaseUrl == "" {
		log.Fatalln("Set GO_DATABASE_URL enviroment variable")
	}
	// Инициализация БД
	db, err := repository.NewPostgres(databaseUrl)
	if err != nil {
		log.Println(err.Error())
	}
	timezone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalln("Error setting time zone")
	}
	goCron := gocron.NewScheduler(timezone)
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	bot := telegram.NewBot(botApi, services, goCron, adminsList)
	if err := bot.StartJobs(); err != nil {
		log.Fatal(err)
	}
	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}
