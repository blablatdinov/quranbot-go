package service

import "github.com/blablatdinov/quranbot-go/internal/storage"

type Bot interface {
	GetOrCreateSubscriber(chatId int64) (string, error)
	RegisterSubscriber(chatId int64) error
}

type Service struct {
	Bot
}

func NewService(repos *storage.Repository) *Service {
	return &Service{
		Bot: NewBotService(repos.Bot),
	}
}
