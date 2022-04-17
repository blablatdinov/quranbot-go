package main

import (
	"fmt"
	"github.com/blablatdinov/quranbot-go/internal/service"
	"github.com/blablatdinov/quranbot-go/internal/transport/rest"
	"github.com/blablatdinov/quranbot-go/internal/transport/rest/handler"
	"log"
	"os"

	"github.com/blablatdinov/quranbot-go/internal/storage"
	_ "github.com/lib/pq"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalf("error loading env variables: %s", err.Error())
	//}
	//bot := telegramsdk.NewBot(os.Getenv("BOT_TOKEN"))
	// updatesChan := bot.GetUpdatesChan()
	//updatesChan, err := bot.RunWebhookServer(os.Getenv("HOST"))
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	databaseUrl := os.Getenv("GO_DATABASE_URL")
	fmt.Println("Database url: ", databaseUrl)
	db, err := storage.NewPostgres(databaseUrl)
	repos := storage.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	if err != nil {
		log.Fatal(err.Error())
	}
	server := new(rest.Server)
	if err := server.Run("8001", handlers.InitRoutes()); err != nil {
		 log.Fatalf("Rest server error: %s", err.Error())
	}

	//log.Println("Rest server started")
	//
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	//
	//log.Println("Rest shutting down")
	//
	//if err := db.Close(); err != nil {
	//	log.Fatalln(err.Error())
	//}

	//for message := range updatesChan {
	//	fmt.Printf("message text: %s\n", message.Text)
	//	if len(message.Text) > 5 && message.Text[:6] == "/start" {
	//		referralCode := "0"
	//		if len(message.Text) > 6 {
	//			referralCode = message.Text[7:]
	//		}
	//		text, err := services.GetOrCreateSubscriber(message.Chat.Id, referralCode)
	//		if err != nil {
	//			log.Fatal(err.Error())
	//		}
	//		bot.SendMessage(message.Chat.Id, text)
	//	}
	//}
}
