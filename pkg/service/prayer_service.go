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

// GetPrayer получить текст и клавиатуру с временем намаза
func (s *PrayerService) GetPrayer(chatId int64, targetTime time.Time) (string, tgbotapi.InlineKeyboardMarkup, error) {
	subscriberHasCity, err := s.repo.SubscriberHasCity(chatId)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	if !subscriberHasCity {
		return "", tgbotapi.InlineKeyboardMarkup{}, errors.New("subscriber hasn't city")
	}
	prayers, err := s.repo.GetPrayer(chatId, targetTime)
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

func (s *PrayerService) GetCityByName(cityName string) (qbot.City, error) {
	city, err := s.repo.GetCityByName(cityName)
	if err != nil {
		return qbot.City{}, err
	}
	return city, nil
}

func (s *PrayerService) ChangeCity(chatId int64, cityId int) error {
	err := s.repo.ChangeCity(chatId, cityId)
	return err
}

// GetPrayersForMailing получить время намаза для пользователей
// Каждый день в 20:00 проходит рассылка с временами намаза для след. дня
// TODO: сейчас на каждого подписчика выполняется отдельный запрос в БД, оптимизировать
func (s *PrayerService) GetPrayersForMailing() ([]qbot.Answer, error) {
	subscriberWithCityChatIds, err := s.repo.GetSubscriberWithCityChatIds()
	var result []qbot.Answer
	if err != nil {
		return []qbot.Answer{}, err
	}
	for _, chatId := range subscriberWithCityChatIds {
		content, keyboard, err := s.GetPrayer(chatId, time.Now().AddDate(0, 0, 1))
		result = append(result, qbot.Answer{ChatId: chatId, Content: content, Keyboard: keyboard})
		if err != nil {
			return []qbot.Answer{}, err
		}
	}
	return result, nil
}
