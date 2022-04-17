package service

import (
	"github.com/blablatdinov/quranbot-go/internal/core"
	"github.com/blablatdinov/quranbot-go/internal/storage"
	mock_repository "github.com/blablatdinov/quranbot-go/internal/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContentService_GetAyat(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockContent, suraAyat string)
	testTable := []struct {
		name          string
		query         string
		expectedValue int
		mockBehavior  mockBehavior
		expectError   bool
		errorText     string
	}{
		//{
		//	name: "invalid sura num",
		//	query: "a1:1",
		//	expectedValue: "asdf",
		//	mockBehavior: func(s *mock_repository.MockContent, suraAyat string) {
		//		//s.EXPECT().GetSura
		//	},
		//},
		//{
		//	name: "sura number out of range",
		//	query: "1000:1",
		//	expectedValue: "asdf",
		//	mockBehavior: func(s *mock_repository.MockContent, suraAyat string) {
		//		//s.EXPECT().GetSura
		//	},
		//},
		//{
		//	name: "ayat number out of range",
		//	query: "1:8",
		//	expectedValue: "asdf",
		//	mockBehavior: func(s *mock_repository.MockContent, suraAyat string) {
		//		//s.EXPECT().GetSura
		//	},
		//},
		{
			name:          "test number with hypen",
			query:         "1:1",
			expectedValue: 1,
			expectError:   false,
			mockBehavior: func(s *mock_repository.MockContent, suraAyat string) {
				s.EXPECT().GetAyatsBySuraNum(1).Return([]core.Ayat{{Id: 1, Ayat: "1-7"}}, nil)
			},
		},
		{
			name:          "ayat number out of range",
			query:         "1:8",
			expectedValue: 0,
			expectError:   true,
			errorText:     "ayat index out of range",
			mockBehavior: func(s *mock_repository.MockContent, suraAyat string) {
				s.EXPECT().GetAyatsBySuraNum(1).Return([]core.Ayat{{Id: 1, Ayat: "1-7"}}, nil)
			},
		},
		{
			name:          "sura number out of range",
			query:         "1000:1",
			expectedValue: 0,
			expectError:   true,
			errorText:     "sura index out of range",
			mockBehavior:  func(s *mock_repository.MockContent, suraAyat string) {},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			contentRepository := mock_repository.NewMockContent(c)
			testCase.mockBehavior(contentRepository, testCase.query)
			repository := &storage.Repository{Content: contentRepository}
			service := NewService(repository)

			answer, err := service.GetAyatBySuraAyatNum(testCase.query)
			if testCase.expectError {
				assert.Equal(t, testCase.errorText, err.Error())
			}
			assert.Equal(t, testCase.expectedValue, answer)
		})
	}
}
