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

func getPrayerAtUserKeyboard(prayers []qbot.PrayerAtUser) tgbotapi.InlineKeyboardMarkup {
	buttons := make([]tgbotapi.InlineKeyboardButton, 0, 5)
	for _, prayerAtUser := range prayers {
		var buttonEmoji string
		var buttonData string
		if prayerAtUser.IsRead {
			buttonEmoji = "✅"
			buttonData = fmt.Sprintf("setPrayerStatusToUnread(%d)", prayerAtUser.Id)
		} else {
			buttonEmoji = "❌"
			buttonData = fmt.Sprintf("setPrayerStatusToRead(%d)", prayerAtUser.Id)
		}
		buttons = append(
			buttons,
			tgbotapi.NewInlineKeyboardButtonData(
				buttonEmoji,
				buttonData,
			),
		)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	return keyboard
}

func (s *PrayerService) getKeyboardWithPrayers(chatId int64, prayers []qbot.Prayer) (tgbotapi.InlineKeyboardMarkup, error) {
	prayersAtUser, err := s.repo.GetOrCreatePrayerForUser(chatId, prayers)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	return getPrayerAtUserKeyboard(prayersAtUser), nil
}

func (s *PrayerService) ChangePrayerStatus(prayerAtUserId int, status bool) (tgbotapi.InlineKeyboardMarkup, error) {
	err := s.repo.ChangePrayerStatus(prayerAtUserId, status)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	prayersAtUser, err := s.repo.GetPrayersAtUserByOnePrayerId(prayerAtUserId)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	keyboard := getPrayerAtUserKeyboard(prayersAtUser)
	return keyboard, err
}
