package service

import (
	"errors"
	"github.com/blablatdinov/quranbot-go/internal/core"
	"github.com/blablatdinov/quranbot-go/internal/storage"
	"strconv"
	"strings"
)

type ContentService struct {
	repo storage.Content
}

func NewContentService(repo storage.Content) *ContentService {
	return &ContentService{repo}
}

func (s *ContentService) GetAyatBySuraAyatNum(query string) (int, error) {
	splittedQuery := strings.Split(query, ":")
	suraNum, err := strconv.Atoi(strings.TrimSpace(splittedQuery[0]))
	if err != nil {
		return 0, err
	}
	if suraNum < 1 || suraNum > 114 {
		return 0, errors.New("sura index out of range")
	}
	ayats, err := s.repo.GetAyatsBySuraNum(suraNum)
	if err != nil {
		return 0, err
	}
	for _, ayat := range ayats {
		//fmt.Printf("Iterate by ayats %d\n", ayat.Id)
		if checkAyatInResult(splittedQuery[1], ayat) {
			return ayat.Id, nil
		}
	}
	return 0, errors.New("ayat index out of range")
}

func checkAyatInResult(query string, ayat core.Ayat) bool {
	//fmt.Printf("checkAyatInResult: %s, %s\n", query, ayat.Ayat)
	switch {
	case strings.Contains(ayat.Ayat, "-"):
		return serviceAyatRange(query, ayat.Ayat, "-")
	case strings.Contains(ayat.Ayat, ","):
		return serviceAyatRange(query, ayat.Ayat, ",")
	case query == ayat.Ayat:
		return true
	default:
		return false
	}
}

func serviceAyatRange(query, ayat, separator string) bool {
	separatedAyatNum := strings.Split(ayat, separator)
	leftLimit, err := strconv.Atoi(strings.TrimSpace(separatedAyatNum[0]))
	if err != nil {
		return false
	}
	rightLimit, err := strconv.Atoi(strings.TrimSpace(separatedAyatNum[1]))
	if err != nil {
		return false
	}
	queryAyatInt, err := strconv.Atoi(query)
	if err != nil {
		return false
	}
	for leftLimit <= rightLimit {
		if queryAyatInt == leftLimit {
			return true
		}
		leftLimit++
	}
	return false
}
