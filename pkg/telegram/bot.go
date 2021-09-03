package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"qbot/pkg/service"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
}

func NewBot(bot *tgbotapi.BotAPI, service *service.Service) *Bot {
	return &Bot{
		bot: bot,
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