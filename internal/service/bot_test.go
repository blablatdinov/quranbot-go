package service

import (
	"errors"
	"github.com/blablatdinov/quranbot-go/internal/core"
	"github.com/blablatdinov/quranbot-go/internal/storage"
	mock_repository "github.com/blablatdinov/quranbot-go/internal/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBotService_GetOrCreateSubscriber(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockBot, chatId int64, referralCode string)
	testTable := []struct {
		name          string
		chatId        int64
		referralCode  string
		mockBehavior  mockBehavior
		expectedValue string
	}{
		{
			name:          "Test new user",
			chatId:        29833948,
			referralCode:  "",
			expectedValue: "register",
			mockBehavior: func(s *mock_repository.MockBot, chatId int64, referralCode string) {
				s.EXPECT().GetSubscriberByChatId(chatId).Return(core.Subscriber{IsActive: false}, errors.New("sql: no rows in result set"))
			},
		},
		{
			name:          "Test reactivate user",
			chatId:        29833948,
			referralCode:  "",
			expectedValue: "Рады видеть вас снова, вы продолжите с дня 3",
			mockBehavior: func(s *mock_repository.MockBot, chatId int64, referralCode string) {
				s.EXPECT().GetSubscriberByChatId(chatId).Return(core.Subscriber{Day: 3}, nil)
			},
		},
		{
			name:          "Test for already active user",
			chatId:        29833948,
			referralCode:  "",
			expectedValue: "Вы уже зарегистрированы",
			mockBehavior: func(s *mock_repository.MockBot, chatId int64, referralCode string) {
				s.EXPECT().GetSubscriberByChatId(chatId).Return(core.Subscriber{IsActive: true}, nil)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			botRepository := mock_repository.NewMockBot(c)
			testCase.mockBehavior(botRepository, testCase.chatId, testCase.referralCode)
			repository := &storage.Repository{Bot: botRepository}
			service := NewService(repository)
			answer, _ := service.GetOrCreateSubscriber(testCase.chatId, testCase.referralCode)
			assert.Equal(t, testCase.expectedValue, answer)
		})
	}
}
