package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"qbot"
	"qbot/pkg/service"
	"sync"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
	goCron  *gocron.Scheduler
}

func NewBot(bot *tgbotapi.BotAPI, service *service.Service, goCron *gocron.Scheduler) *Bot {
	return &Bot{
		bot:     bot,
		service: service,
		goCron:  goCron,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	var updates tgbotapi.UpdatesChannel
	webhookHost := os.Getenv("WEBHOOK_HOST")

	if webhookHost == "" {
		_, err := b.bot.RemoveWebhook()
		if err != nil {
			return err
		}
		updates, err = b.bot.GetUpdatesChan(u)
		if err != nil {
			return err
		}
		log.Println("Bot started on long polling...")
	} else {
		_, err := b.bot.RemoveWebhook()
		if err != nil {
			return err
		}
		log.Printf("Setting webhook on %s...\n", webhookHost)
		_, err = b.bot.SetWebhook(tgbotapi.NewWebhook(webhookHost + b.bot.Token))
		if err != nil {
			log.Fatalf("Setting webhook error: %s", err.Error())
			return err
		}
		log.Println("Getting webhook info...")
		info, err := b.bot.GetWebhookInfo()
		if err != nil {
			return err
		}
		if info.LastErrorDate != 0 {
			return errors.New(fmt.Sprintf("Telegram callback failed: %s", info.LastErrorMessage))
		}
		log.Println("Getting updates channel...")
		updates = b.bot.ListenForWebhook("/" + b.bot.Token)
		go http.ListenAndServe("localhost:8012", nil)
		log.Println("Bot started on webhook...")
	}

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			b.handleError(0, errors.New("unknown behaviour"))
			continue
		}
		if update.CallbackQuery != nil {
			if err := b.handleQuery(update.CallbackQuery); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}
		if update.Message.IsCommand() {
			log.Printf("Take command\n")
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
	return nil
}

func (b *Bot) ReadMessagesChan(quitFromReadLoop chan struct{}, messagesChan chan tgbotapi.Message, messageListChan chan []tgbotapi.Message) {
	var messages []tgbotapi.Message
	for {
		select {
		case message := <-messagesChan:
			messages = append(messages, message)
		case <-quitFromReadLoop:
			goto ENDLOOP
		default:
			continue
		}
	}
ENDLOOP:
	messageListChan <- messages
}

// MassMailing функция для массовой рассылки
func (b *Bot) MassMailing(content []qbot.Answer) ([]int64, error) {
	messagesChan := make(chan tgbotapi.Message, len(content))
	var wg sync.WaitGroup
	for _, elem := range content {
		wg.Add(1)
		_elem := elem
		go func(messagesChan chan tgbotapi.Message, wg *sync.WaitGroup) {
			log.Printf("Send mailing to %d (%s)", _elem.ChatId, _elem.Content[:50])
			message, err := b.SendMessageV2(_elem)
			if err != nil {
				log.Printf("Error: %s", err.Error())
				wg.Done()
				return
			}
			messagesChan <- message
			wg.Done()
		}(messagesChan, &wg)
	}
	quitFromReadLoop := make(chan struct{})
	messageListChan := make(chan []tgbotapi.Message, 1)
	go b.ReadMessagesChan(quitFromReadLoop, messagesChan, messageListChan)
	wg.Wait()
	quitFromReadLoop <- struct{}{}
	sendedMessages := <-messageListChan
	var sendedToChatIds []int64
	for _, c := range sendedMessages {
		sendedToChatIds = append(sendedToChatIds, c.Chat.ID)
	}
	deactivatedSubscriptionIds := difference(sendedMessages, content)
	if len(deactivatedSubscriptionIds) > 0 {
		if err := b.service.DeactivateSubscribers(deactivatedSubscriptionIds); err != nil {
			return []int64{}, err
		}
	}
	return sendedToChatIds, nil
}

// SendMorningContent рассылка аятов
// каждое утро в 7:00
func (b *Bot) SendMorningContent() error {
	log.Println("Send morning content task started...")
	content, err := b.service.GetMorningContentForTodayMailing()
	if err != nil {
		return err
	}
	log.Println("Content length:", len(content))
	chatIdsForUpdateDay, err := b.MassMailing(content)
	if err != nil {
		return err
	}
	err = b.service.UpdateDaysForSubscribers(chatIdsForUpdateDay)
	if err != nil {
		return err
	}
	return err
}

// SendPrayerTimes Рассылка времени намаза для след. дня
func (b *Bot) SendPrayerTimes() error {
	log.Println("Send prayer times task started...")
	prayerTimesAtUser, err := b.service.GetPrayersForMailing()
	if err != nil {
		return err
	}
	if _, err = b.MassMailing(prayerTimesAtUser); err != nil {
		return err
	}
	return nil
}

// difference Найти разницу в двух массивах по ChatId
// Кейс: произошла рассылка, некоторые подписчики заблокировали бота,
// чтобы найти тех, кто отписался находим разницу между тем, что должно было отправиться и что фактически отправилось
func difference(messages []tgbotapi.Message, contents []qbot.Answer) []int64 {
	mb := make(map[int64]struct{}, len(messages))
	for _, x := range messages {
		mb[x.Chat.ID] = struct{}{}
	}
	var diff []int64
	for _, x := range contents {
		if _, found := mb[x.ChatId]; !found {
			diff = append(diff, x.ChatId)
		}
	}
	return diff
}
