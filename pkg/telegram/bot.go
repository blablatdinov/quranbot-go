package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"qbot"
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

func (b *Bot) SendMorningContent() error {
	content, err := b.service.GetMorningContentForTodayMailing()
	messagesChan := make(chan tgbotapi.Message, len(content))
	var wg sync.WaitGroup
	if err != nil {
		return err
	}
	for _, elem := range content {
		wg.Add(1)
		chatId := elem.ChatId
		content := elem.Content
		go func(messagesChan chan tgbotapi.Message, wg *sync.WaitGroup) {
			message, err := b.SendMessage(chatId, content)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
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
	messages := <-messageListChan
	inactivatedSubscribers := difference(messages, content)
	err = b.service.DeactivateSubscribers(inactivatedSubscribers)
	return err
}

func difference(messages []tgbotapi.Message, contents []qbot.MailingContent) []int64 {
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
