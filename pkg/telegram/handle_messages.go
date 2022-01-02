package telegram

import (
	"errors"
	"fmt"
	"log"
	"qbot"
	"qbot/pkg/service"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("handleMessage: Message \"%s\" from %d\n", message.Text, message.Chat.ID)
	patterns := map[string]interface{}{
		`/start`:        b.handleStartCommand,
		`\d.?:.?\d`:     b.searchAyatBySuraAyatNum,
		`(И|и)збранное`: b.getFavoriteAyats,
		`(П|п)одкасты`:  b.getRandomPodcast,
		`Время намаза`:  b.getPrayerTimes,
		`.+`:            b.changeCity,
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

func (b Bot) changeCity(message *tgbotapi.Message) error {
	city, err := b.service.GetCityByName(message.Text)
	if err != nil {
		return err
	}
	if city.Id == 0 {
		return errors.New("city not found")
	}
	err = b.service.ChangeCity(message.Chat.ID, city.Id)
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Вам будет приходить время намаза для г. %s", city.Name))
	b.bot.Send(msg)
	return err
}

// getPrayerTimes получить время намазов для пользователя
func (b Bot) getPrayerTimes(message *tgbotapi.Message) error {
	answer, keyboard, err := b.service.GetPrayer(message.Chat.ID, time.Now())
	if err != nil {
		if err.Error() == "subscriber hasn't city" {
			answer = "subscriber hasn't city"
		} else {
			return err
		}
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, answer)
	msg.ReplyMarkup = keyboard
	b.bot.Send(msg)
	return nil
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
		b.SendMessage(
			qbot.Answer{
				ChatId:   message.Chat.ID,
				Content:  podcast.LinkToFile,
				Keyboard: tgbotapi.InlineKeyboardMarkup{},
			},
		)
	}
	if err != nil {
		return err
	}
	return nil
}

func (b Bot) getFavoriteAyats(message *tgbotapi.Message) error {
	answer, keyboard, err := b.service.GetFavoriteAyats(message.Chat.ID)
	if err != nil {
		if err.Error() == "subscriber hasn't favorite ayats" {
			answer = "Вы не добавили аятов в избранное"
			b.SendMessage(
				qbot.Answer{
					ChatId:   message.Chat.ID,
					Content:  answer,
					Keyboard: tgbotapi.InlineKeyboardMarkup{},
				},
			)
			return nil
		} else {
			return err
		}
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
		msg.ReplyMarkup = service.GetDefaultKeyboard()
	} else if err != nil && err.Error() == "ayat not found" {
		msg.Text = ayatNotFoundText
		msg.ReplyMarkup = service.GetDefaultKeyboard()
	} else if err != nil {
		return err
	}
	log.Printf("Exit from if\n")
	msg.ParseMode = "markdown"
	_, err = b.bot.Send(msg)
	return err
}
