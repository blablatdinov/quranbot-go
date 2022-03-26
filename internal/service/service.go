package service

import "github.com/blablatdinov/quranbot-go/internal/storage"

type Bot interface {
	GetOrCreateSubscriber(chatId int64, referralCode string) (string, error)
	RegisterSubscriber(chatId int64, referralCode string) error
}

type Service struct {
	Bot
}

func NewService(repos *storage.Repository) *Service {
	return &Service{
		Bot: NewBotService(repos.Bot),
	}
}
