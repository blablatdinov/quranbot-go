package telegram

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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
		if err.Error() != "len(chatIds) must be more 0" {
			return err
		}
	}
	log.Printf("Unsubscribed %d users\n", len(deactivatedSubscribersIds))
	activeSubscribersCount, err := b.service.GetSubscribersCount("active")
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	totalSubscribersCount, err := b.service.GetSubscribersCount("total")
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	log.Printf("Active users count: %d\n", activeSubscribersCount)
	log.Printf("Total users count: %d\n", totalSubscribersCount)
	return nil
}
