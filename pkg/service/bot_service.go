package service

import (
	"fmt"
	"log"
	"qbot/pkg/repository"
)

type BotService struct {
	repo repository.Bot
}

func NewBotService(repo repository.Bot) *BotService {
	return &BotService{
		repo,
	}
}

func (s *BotService) GetAyatByMailingDay(mailingDay int) (string, error) {
	db, err := repository.NewPostgres()
	if err != nil {
		log.Fatal(err.Error())
	}
	contentService := NewContentService(repository.NewContentPostgres(db))
	content, err := contentService.GetAyatByMailingDay(mailingDay)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (s *BotService) CreateSubscriber(chatId int64) (string, bool) {
	subscriber, created, err := s.repo.GetOrCreateSubscriber(chatId)
	if err != nil {
		log.Fatal(err)
	}
	if created {
		if err != nil {
			log.Fatal(err)
		}
		return "Вы успешно зарегестрировались", created
	} else {
		if subscriber.IsActive {
			return "Вы уже зарегестрированы", created
		} else {
			return fmt.Sprintf("Рады видеть вас снова, вы продолжите с дня %d", subscriber.Day), created
		}
	}
}

func (s *BotService)SetSubscriberState(chatId int64, step string) error{
	err := s.repo.SetSubscriberState(chatId, step)
	return err
}
