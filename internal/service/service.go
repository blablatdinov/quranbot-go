package service

import "github.com/blablatdinov/quranbot-go/internal/storage"

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Bot interface {
	GetOrCreateSubscriber(chatId int64, referralCode string) (string, error)
	RegisterSubscriber(chatId int64, referralCode string) error
}

type Content interface {
	GetAyatBySuraAyatNum(suraAyat string) (int, error)
}

type Service struct {
	Bot
	Content
}

func NewService(repos *storage.Repository) *Service {
	return &Service{
		Bot:     NewBotService(repos.Bot),
		Content: NewContentService(repos.Content),
	}
}
