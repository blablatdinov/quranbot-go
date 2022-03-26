package service

import (
	"fmt"

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
	fmt.Println(subscriber.IsActive)
	if err == nil {
		if subscriber.IsActive == true {
			return "Вы уже зарегистрированы", nil
		} else if subscriber.IsActive == false {
			s.ActivateSubscriber(chatId)
			return fmt.Sprintf("Рады видеть вас снова, вы продолжите с дня %d", subscriber.Day), nil
		} else if err.Error() == "sql: no rows in result set" {
			if err := s.RegisterSubscriber(chatId, referralCode); err != nil {
				return "", err
			}
			return "register", nil
		}
	}
	return "", err
}

func (s *BotService) RegisterSubscriber(chatId int64, referralCode string) error {
	return s.repo.CreateSubscriber(chatId, referralCode)
}

func (s *BotService) ActivateSubscriber(chatId int64) error {
	return s.repo.ActivateSubscriber(chatId)
}
