package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"qbot/pkg/service"
	"sync"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
}

func NewBot(bot *tgbotapi.BotAPI, service *service.Service) *Bot {
	return &Bot{
		bot:     bot,
		service: service,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			b.handleQuery(update.CallbackQuery)
			continue
		}
		if update.Message.IsCommand() {
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

func (b *Bot) CheckSubscriberStatus(chatId int64, deactivatedSubscribersChan chan int64, limit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	limit <- struct{}{}
	action := tgbotapi.NewChatAction(chatId, "typing")
	_, err := b.bot.Send(action)
	if err != nil {
		if err.Error() == "Bad Request: chat not found" || err.Error() == "Forbidden: bot was blocked by the user" {
			deactivatedSubscribersChan <- chatId
		}
	}
	<-limit
}

func (b *Bot) SomeFunc(deactivatedSubscribersChan chan int64, quitFromCycle chan struct{}, deactivatedSubscriberIdsChan chan []int64) { // TODO: rename
	var deactivatedSubscriberIds []int64
	for {
		select {
		case deactivatedSubscriberId := <-deactivatedSubscribersChan:
			deactivatedSubscriberIds = append(deactivatedSubscriberIds, deactivatedSubscriberId)
		case <-quitFromCycle:
			goto BREAKLOOP
		default:
			continue
		}
	}
BREAKLOOP:
	deactivatedSubscriberIdsChan <- deactivatedSubscriberIds
}

func (b *Bot) CheckSubscribers() error {
	subsribers, err := b.service.GetActiveSubscribers()
	deactivatedSubscribersChan := make(chan int64)
	deactivatedSubscribersIdsChan := make(chan []int64)
	quitFromCycle := make(chan struct{}, 1)
	var wg sync.WaitGroup
	limit := make(chan struct{}, 100)
	if err != nil {
		return err
	}
	go b.SomeFunc(deactivatedSubscribersChan, quitFromCycle, deactivatedSubscribersIdsChan)
	for _, subsriber := range subsribers {
		wg.Add(1)
		go b.CheckSubscriberStatus(subsriber.ChatId, deactivatedSubscribersChan, limit, &wg)
	}
	wg.Wait()
	quitFromCycle <- struct{}{}
	deactivatedSubscribersIds := <-deactivatedSubscribersIdsChan
	err = b.service.DeactivateSubscribers(deactivatedSubscribersIds)
	if err != nil {
		if err.Error() == "len(chatIds) must be more 0" {
			log.Println("Unsubscribed 0 users")
			return nil
		}
		return err
	}
	log.Printf("Unsubscribed %d users\n", len(deactivatedSubscribersIds))
	return nil
}
