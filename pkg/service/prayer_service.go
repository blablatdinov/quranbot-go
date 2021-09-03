package service

import (
	"errors"
	"log"
	"qbot/pkg/repository"
)

type PrayerService struct {
	repo repository.Prayer
}

func NewPrayerService(repo repository.Prayer) *PrayerService {
	return &PrayerService{repo}
}

func (s *PrayerService) GetPrayer(chatId int64) (string, error) {
	subscriberHasCity, err := s.repo.SubscriberHasCity(chatId)
	if err != nil {
		return "", err
	}
	if !subscriberHasCity {
		return "", errors.New("subscriber hasn't city")
	}
	prayers, err := s.repo.GetPrayer(chatId)
	log.Printf("%d\n", len(prayers))
	if len(prayers) == 0 {
		return "", errors.New("null prayers")
	}
	return prayers[0].Name, err
}
