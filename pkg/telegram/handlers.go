package telegram

import (
	"log"
	"qbot"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const commandStart = "start"

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	log.Printf("Command: %s\n", message.Command())
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	answers := b.service.CreateSubscriber(message.Chat.ID)
	for _, answer := range answers {
		b.SendMessage(answer)
	}
	return nil
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	// register logic
	return nil
}

func (b *Bot) handleError(chatId int64, err error) error {
	log.Printf("handleError: %s\n", err.Error())
	return nil
}

func (b *Bot) SendMessage(answer qbot.Answer) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(answer.ChatId, answer.Content)
	msg.ParseMode = "markdown"
	if answer.HasKeyboard() {
		msg.ReplyMarkup = answer.Keyboard
	}
	message, err := b.bot.Send(msg)
	return message, err
}

func path(expression string, message string) bool {
	matched, err := regexp.MatchString(expression, message)
	if err != nil {
		return false
	}
	return matched
}
