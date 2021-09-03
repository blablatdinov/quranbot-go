package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("handleMessage: Message \"%s\" from %d\n", message.Text, message.Chat.ID)
	patterns := map[string]interface{}{
		`\d.?:.?\d`:     b.searchAyatBySuraAyatNum,
		`(И|и)збранное`: b.getFavoriteAyats,
		`(П|п)одкасты`:  b.getRandomPodcast,
	}
	for pattern, handler := range patterns {
		if path(pattern, message.Text) {
			err := handler.(func(*tgbotapi.Message) error)(message)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("unknow pattern")
}

func (b Bot) getRandomPodcast(message *tgbotapi.Message) error {
	err := b.service.SetSubscriberState(message.Chat.ID, "")
	if err != nil {
		return err
	}
	podcast, err := b.service.GetRandomPodcast()
	if podcast.TgFileId.Valid {
		msg := tgbotapi.NewAudioShare(message.Chat.ID, podcast.TgFileId.String)
		_, err := b.bot.Send(msg)
		if err != nil {
			log.Printf("handler: %s", err.Error())
		}
	} else {
		b.SendMessage(message.Chat.ID, podcast.LinkToFile)
	}
	if err != nil {
		return err
	}
	return nil
}

func (b Bot) getFavoriteAyats(message *tgbotapi.Message) error {
	answer, keyboard, err := b.service.GetFavoriteAyats(message.Chat.ID)
	if err != nil {
		return err
	}
	err = b.service.SetSubscriberState(message.Chat.ID, "see favorite")
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, answer)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = keyboard
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (b Bot) searchAyatBySuraAyatNum(message *tgbotapi.Message) error {
	err := b.service.SetSubscriberState(message.Chat.ID, "")
	if err != nil {
		return err
	}
	log.Printf("searchAyatBySuraAyatNum: search '%s' ayat\n", message.Text)
	answer, keyboard, err := b.service.GetAyatBySuraAyatNum(message.Chat.ID, message.Text, "")
	ayatNotFoundText := "Аят не найден"
	suraNotFoundText := "Сура не найдена"
	msg := tgbotapi.NewMessage(message.Chat.ID, answer)
	msg.ReplyMarkup = keyboard
	if err != nil && err.Error() == "sura not found" {
		msg.Text = suraNotFoundText
		msg.ReplyMarkup = b.service.GetDefaultKeyboard()
	} else if err != nil && err.Error() == "ayat not found" {
		msg.Text = ayatNotFoundText
		msg.ReplyMarkup = b.service.GetDefaultKeyboard()
	} else if err != nil {
		return err
	}
	log.Printf("Exit from if\n")
	msg.ParseMode = "markdown"
	_, err = b.bot.Send(msg)
	return err
}
