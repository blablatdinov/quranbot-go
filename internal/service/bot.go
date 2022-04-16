package service

import (
	"fmt"
	"github.com/blablatdinov/quranbot-go/internal/core"
	"github.com/blablatdinov/quranbot-go/internal/storage"
)

type BotService struct {
	repo storage.Bot
}

func NewBotService(repo storage.Bot) *BotService {
	return &BotService{
		repo,
	}
}

func (s *BotService) GetOrCreateSubscriber(chatId int64, referralCode string) (string, error) {
	subscriber, err := s.repo.GetSubscriberByChatId(chatId)
	if isNewUser(err, subscriber) {
		return "register", nil
	} else if isReactivatedUser(err, subscriber) {
		return fmt.Sprintf("Рады видеть вас снова, вы продолжите с дня %d", subscriber.Day), nil
	} else if isAlreadyActiveUser(err, subscriber) {
		return "Вы уже зарегистрированы", nil
	}
	return "", err
}

func isNewUser(err error, subscriber core.Subscriber) bool {
	return err != nil && err.Error() == "sql: no rows in result set" && subscriber.IsActive == false
}

func isReactivatedUser(err error, subscriber core.Subscriber) bool {
	return err == nil && subscriber.IsActive == false
}

func isAlreadyActiveUser(err error, subscriber core.Subscriber) bool {
	return err == nil && subscriber.IsActive == true
}

func (s *BotService) RegisterSubscriber(chatId int64, referralCode string) error {
	return s.repo.CreateSubscriber(chatId, referralCode)
}

func (s *BotService) ActivateSubscriber(chatId int64) error {
	return s.repo.ActivateSubscriber(chatId)
}
