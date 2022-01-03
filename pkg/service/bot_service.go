package service

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"qbot"
	"qbot/pkg/repository"
	"time"
)

type BotService struct {
	repo repository.Bot
}

func NewBotService(repo repository.Bot) *BotService {
	return &BotService{
		repo,
	}
}

func (s *BotService) CreateSubscriber(chatId int64) []qbot.Answer {
	subscriber, created, err := s.repo.GetOrCreateSubscriber(chatId)
	if err != nil {
		log.Fatal(err)
	}
	if created {
		if err != nil {
			log.Fatal(err)
		}
		return []qbot.Answer{{
			ChatId:   chatId,
			Content:  "Вы успешно зарегестрировались",
			Keyboard: tgbotapi.InlineKeyboardMarkup{},
		}}
	} else {
		if subscriber.IsActive {
			return []qbot.Answer{{
				ChatId:   chatId,
				Content:  "Вы уже зарегестрированы",
				Keyboard: tgbotapi.InlineKeyboardMarkup{},
			}}
		} else {
			return []qbot.Answer{{
				ChatId:   chatId,
				Content:  fmt.Sprintf("Рады видеть вас снова, вы продолжите с дня %d", subscriber.Day),
				Keyboard: tgbotapi.InlineKeyboardMarkup{},
			}}
		}
	}
}

func (s *BotService) SetSubscriberState(chatId int64, step string) error {
	err := s.repo.SetSubscriberState(chatId, step)
	return err
}

func (s *BotService) GetSubscriberState(chatId int64) (string, error) {
	state, err := s.repo.GetSubscriberState(chatId)
	if err != nil {
		return "", err
	}
	return state, err
}

func (s *BotService) GetAyatByMailingDay(mailingDay int) (string, error) {
	ayat, err := s.repo.GetAyatByMailingDay(mailingDay)
	contentTemplate := "%d: %s) %s\n\nСсылка на %s"
	suraLink := fmt.Sprintf("[источник](https://umma.ru%s)", ayat.SuraLink)
	content := fmt.Sprintf(contentTemplate, 1, ayat.Ayat, ayat.Content, suraLink)
	if err != nil {
		return "", err
	}
	return content, err
}

func (s *BotService) GetActiveSubscribers() ([]qbot.Subscriber, error) {
	subscribers, err := s.repo.GetActiveSubscribers()
	return subscribers, err
}

func (s *BotService) DeactivateSubscribers(chatIds []int64) error {
	if len(chatIds) == 0 {
		return errors.New("len(chatIds) must be more 0")
	}
	err := s.repo.DeactivateSubscribers(chatIds)
	return err
}

func (s *BotService) GetSubscribersCount(param string) (int, error) {
	res, err := s.repo.GetSubscribersCount(param)
	return res, err
}

func (s *BotService) SaveMessage(message *tgbotapi.Message, isUnknown bool) {
	message_ := s.tgbotapiMessageToQbotMessage(*message, 0, isUnknown)
	if err := s.repo.SaveMessage(message_); err != nil {

	}
}

func (s *BotService) BulkSaveMessages(messages []tgbotapi.Message, mailingId int) {
	messages_ := []qbot.Message{}
	for _, message := range messages {
		messages_ = append(messages_, s.tgbotapiMessageToQbotMessage(message, mailingId, false))
	}
	s.repo.BulkSaveMessages(messages_)
}

func (s *BotService) CreateMailing() (int, error) {
	return s.repo.CreateMailing()
}

func (s *BotService) tgbotapiMessageToQbotMessage(message tgbotapi.Message, mailingId int, isUnknown bool) qbot.Message {
	messageJson, err := json.Marshal(message)
	if err != nil {

	}
	var isUnknown_ string
	if isUnknown {
		isUnknown_ = "t"
	} else {
		isUnknown_ = "f"
	}
	messageJsonString := string(messageJson)
	dateTime := time.Unix(int64(message.Date), 0)
	return qbot.Message{
		Date:       dateTime.Format("2006-02-01T15:04:05") + "+03:00",
		FromUserId: int64(message.From.ID),
		MessageId:  message.MessageID,
		ChatId:     message.Chat.ID,
		Text:       message.Text,
		Json:       messageJsonString,
		Mailing:    mailingId,
		IsUnknown:  isUnknown_,
	}
}
