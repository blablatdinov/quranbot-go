package service

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"qbot"
	"qbot/pkg/repository"
	"time"
)

type PrayerService struct {
	repo repository.Prayer
}

func NewPrayerService(repo repository.Prayer) *PrayerService {
	return &PrayerService{repo}
}

func (s *PrayerService) GetPrayer(chatId int64) (string, tgbotapi.InlineKeyboardMarkup, error) {
	subscriberHasCity, err := s.repo.SubscriberHasCity(chatId)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	if !subscriberHasCity {
		return "", tgbotapi.InlineKeyboardMarkup{}, errors.New("subscriber hasn't city")
	}
	prayers, err := s.repo.GetPrayer(chatId, time.Now())
	if len(prayers) == 0 {
		return "", tgbotapi.InlineKeyboardMarkup{}, errors.New("null prayers")
	}
	messageTemplate := "Время намаза для города: %s (%s)\n\n" +
		"Иртәнге: %s\n" +
		"Восход: %s\n" +
		"Өйлә: %s\n" +
		"Икенде: %s\n" +
		"Ахшам: %s\n" +
		"Ястү: %s\n"
	message := fmt.Sprintf(
		messageTemplate,
		prayers[0].CityName,
		prayers[0].Date.Format("02.01.2006"),
		prayers[0].Time.Format("15:04"),
		prayers[1].Time.Format("15:04"),
		prayers[2].Time.Format("15:04"),
		prayers[3].Time.Format("15:04"),
		prayers[4].Time.Format("15:04"),
		prayers[5].Time.Format("15:04"),
	)
	keyboard, err := s.getKeyboardWithPrayers(chatId, prayers)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	return message, keyboard, err
}

func (s *PrayerService) getKeyboardWithPrayers(chatId int64, prayers []qbot.Prayer) (tgbotapi.InlineKeyboardMarkup, error) {
	prayersAtUser, err := s.repo.GetOrCreatePrayerForUser(chatId, prayers)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	buttons := make([]tgbotapi.InlineKeyboardButton, 0, 5)
	for _, prayerAtUser := range prayersAtUser {
		buttons = append(
			buttons,
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d", prayerAtUser.Id),
				"asdf",
			),
		)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	return keyboard, nil
}
