package service

import (
	"errors"
	"fmt"
	"qbot"
	"qbot/pkg/repository"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ContentService struct {
	repo repository.Content
}

func NewContentService(repo repository.Content) *ContentService {
	return &ContentService{repo}
}

func (s *ContentService) GetAyatByMailingDay(mailingDay int) (string, error) {
	ayat, err := s.repo.GetAyatByMailingDay(mailingDay)
	contentTemplate := "%d: %s) %s\n\nСсылка на %s"
	suraLink := fmt.Sprintf("[источник](https://umma.ru%s)", ayat.SuraLink)
	content := fmt.Sprintf(contentTemplate, 1, ayat.Ayat, ayat.Content, suraLink)
	if err != nil {
		return "", err
	}
	return content, err
}

func (s *ContentService) GetAyatBySuraAyatNum(chatId int64, query string, state string) (string, tgbotapi.InlineKeyboardMarkup, error) {
	splittedQuery := strings.Split(query, ":")
	suraNum, err := strconv.Atoi(strings.TrimSpace(splittedQuery[0]))
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	if suraNum < 1 || suraNum > 114 {
		return "", tgbotapi.InlineKeyboardMarkup{}, errors.New("sura not found")
	}
	ayats, err := s.repo.GetAyatsBySuraNum(suraNum)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	var targetAyat qbot.Ayat
	for i, ayat := range ayats {
		if checkAyatInResult(splittedQuery[1], ayat) {
			targetAyat = ayat
			break
		}
		if i == len(ayats)-1 {
			return "", tgbotapi.InlineKeyboardMarkup{}, errors.New("ayat not found")
		}
	}
	targetAyat.IsFavorite = s.repo.AyatIsFavorite(chatId, targetAyat.Id)
	keyboard, err := s.getAyatKeyboard(chatId, targetAyat, state)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	return renderAyat(targetAyat), keyboard, nil
}

func (s *ContentService) GetFavoriteAyats(chatId int64) (string, tgbotapi.InlineKeyboardMarkup, error) {
	ayats, err := s.repo.GetFavoriteAyats(chatId)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	if len(ayats) < 1 {
		return "", tgbotapi.InlineKeyboardMarkup{}, errors.New("subscriber hasn't favorite ayats")
	}
	keyboard := s.getFavoriteAyatsInlineKeyboard(ayats, 0)
	ayat := renderAyat(ayats[0])
	return ayat, keyboard, nil
}

func (s *ContentService) GetFavoriteAyatsFromKeyboard(chatId int64, ayatId int) (string, tgbotapi.InlineKeyboardMarkup, error) {
	ayats, err := s.repo.GetFavoriteAyats(chatId)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	if len(ayats) < 1 {
		return "", tgbotapi.InlineKeyboardMarkup{}, errors.New("subscriber hasn't favorite ayats")
	}
	ayatIndex := getAyatIndex(ayatId, ayats)
	keyboard := s.getFavoriteAyatsInlineKeyboard(ayats, ayatIndex)
	ayat := renderAyat(ayats[ayatIndex])
	return ayat, keyboard, nil
}

func (s *ContentService) GetRandomPodcast() (qbot.Podcast, error) {
	podcast, err := s.repo.GetRandomPodcast()
	if err != nil {
		return qbot.Podcast{}, err
	}
	return podcast, nil
}

func (s *ContentService) GetAyatById(chatId int64, ayatId int, state string) (string, tgbotapi.InlineKeyboardMarkup, error) {
	ayat, err := s.repo.GetAyatById(chatId, ayatId)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	keyboard, err := s.getAyatKeyboard(chatId, ayat, state)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	return renderAyat(ayat), keyboard, nil
}

func (s *ContentService) AddToFavorite(chatId int64, ayatId int, state string) (string, tgbotapi.InlineKeyboardMarkup, error) {
	if err := s.repo.AddToFavorite(chatId, ayatId); err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	_, keyboard, err := s.GetAyatById(chatId, ayatId, state)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	return "Аят добавлен в избранное", keyboard, nil
}

func (s *ContentService) RemoveFromFavorite(chatId int64, ayatId int, state string) (string, tgbotapi.InlineKeyboardMarkup, error) {
	if err := s.repo.RemoveFromFavorite(chatId, ayatId); err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	_, keyboard, err := s.GetAyatById(chatId, ayatId, state)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}
	return "Аят удален из избранных", keyboard, nil
}

func getAyatIndex(ayatId int, ayats []qbot.Ayat) int {
	for i, ayat := range ayats {
		if ayat.Id == ayatId {
			return i
		}
	}
	return 0
}

func (s *ContentService) getAyatKeyboardFromAyatState(chatId int64, ayat qbot.Ayat, addToFavoriteButton tgbotapi.InlineKeyboardButton) (tgbotapi.InlineKeyboardMarkup, error) {
	var keyboard tgbotapi.InlineKeyboardMarkup
	if ayat.Id == 1 {
		nextAyat, err := s.repo.GetAyatById(chatId, ayat.Id+1)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				addToFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s >", nextAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayat.Id+1),
				),
			),
		)
	} else if ayat.Id == 5737 {
		prevAyat, err := s.repo.GetAyatById(chatId, ayat.Id-1)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				addToFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("< %s", prevAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayat.Id-1),
				),
			),
		)
	} else {
		prevAyat, err := s.repo.GetAyatById(chatId, ayat.Id-1)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		nextAyat, err := s.repo.GetAyatById(chatId, ayat.Id+1)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				addToFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("< %s", prevAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayat.Id-1),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s >", nextAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayat.Id+1),
				),
			),
		)
	}
	return keyboard, nil
}

func (s *ContentService) getAdjacentAyatsKeyboard(chatId int64, ayatId int) (tgbotapi.InlineKeyboardMarkup, error) {
	var keyboard tgbotapi.InlineKeyboardMarkup
	ayats, err := s.repo.GetAdjacentAyats(chatId, ayatId)
	if err != nil {
		return keyboard, err
	}
	removeFromFavoriteButton := tgbotapi.NewInlineKeyboardButtonData(
		"Удалить из избранного",
		fmt.Sprintf("removeFromFavorite(%d)", ayatId),
	)
	if ayatId == ayats[0].Id {
		nextAyat := ayats[1]
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				removeFromFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s >", nextAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", nextAyat.Id),
				),
			),
		)
	} else if ayatId == ayats[len(ayats)].Id {
		prevAyat := ayats[len(ayats)]
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				removeFromFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("< %s", prevAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayats[len(ayats)-1].Id),
				),
			),
		)
	} else {
		prevAyat := ayats[len(ayats)-1]
		nextAyat := ayats[ayatId+1]
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				removeFromFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("< %s", prevAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", prevAyat.Id),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s >", nextAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", nextAyat.Id),
				),
			),
		)
	}
	return keyboard, nil
}

func (s *ContentService) getAyatKeyboardFromFavoriteState(chatId int64, ayat qbot.Ayat, addToFavoriteButton tgbotapi.InlineKeyboardButton) (tgbotapi.InlineKeyboardMarkup, error) {
	var keyboard tgbotapi.InlineKeyboardMarkup
	if ayat.Id == 1 {
		nextAyat, err := s.repo.GetAyatById(chatId, ayat.Id+1)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				addToFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s >", nextAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayat.Id+1),
				),
			),
		)
	} else if ayat.Id == 5737 {
		prevAyat, err := s.repo.GetAyatById(chatId, ayat.Id-1)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				addToFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("< %s", prevAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayat.Id-1),
				),
			),
		)
	} else {
		prevAyat, err := s.repo.GetAyatById(chatId, ayat.Id-1)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		nextAyat, err := s.repo.GetAyatById(chatId, ayat.Id+1)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				addToFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("< %s", prevAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayat.Id-1),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s >", nextAyat.GetSuraAyatNum()),
					fmt.Sprintf("getAyat(%d)", ayat.Id+1),
				),
			),
		)
	}
	return keyboard, nil
}

func (s *ContentService) getAyatKeyboard(chatId int64, ayat qbot.Ayat, state string) (tgbotapi.InlineKeyboardMarkup, error) {
	var textForFavorButton string
	var dataForFavorButtonTemplate string
	if ayat.IsFavorite {
		textForFavorButton = "Удалить из избранного"
		dataForFavorButtonTemplate = "removeFromFavorite(%d)"
	} else {
		textForFavorButton = "Добавить в избранное"
		dataForFavorButtonTemplate = "addToFavorite(%d)"
	}
	addToFavoriteButton := tgbotapi.NewInlineKeyboardButtonData(
		textForFavorButton,
		fmt.Sprintf(dataForFavorButtonTemplate, ayat.Id),
	)
	if state == "" {
		keyboard, err := s.getAyatKeyboardFromAyatState(chatId, ayat, addToFavoriteButton)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		return keyboard, nil
	} else {
		keyboard, err := s.getAdjacentAyatsKeyboard(chatId, ayat.Id)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		return keyboard, nil
	}
}

func getTextAndDataForFavoriteButton(isFavorite bool) (string, string) {
	var textForFavorButton string
	var dataForFavorButtonTemplate string
	if isFavorite {
		textForFavorButton = "Удалить из избранного"
		dataForFavorButtonTemplate = "removeFromFavorite(%d)"
	} else {
		textForFavorButton = "Добавить в избранное"
		dataForFavorButtonTemplate = "addToFavorite(%d)"
	}
	return textForFavorButton, dataForFavorButtonTemplate
}

func (s *ContentService) GetMorningContentForTodayMailing() ([]qbot.Answer, error) {
	var result []qbot.Answer
	contentsForSubscriber, err := s.repo.GetMorningContentForTodayMailing()
	if err != nil {
		return []qbot.Answer{}, err
	}
	relativeUrl := strings.Split(contentsForSubscriber[0].Link, "|")[0]
	link := fmt.Sprintf("Ссылка на [источник](https://umma.ru%s)", relativeUrl)
	for _, content := range contentsForSubscriber {
		answer := qbot.Answer{
			ChatId:   content.ChatId,
			Content:  content.Content + link,
			Keyboard: tgbotapi.InlineKeyboardMarkup{},
		}
		fmt.Printf("Service layer: GetMorningContentForTodayMailing: %d %s\n", answer.ChatId, answer.Content)
		result = append(result, answer)
	}
	return result, nil
}

func (s *ContentService) UpdateDaysForSubscribers(chatIds []int64) error {
	if len(chatIds) == 0 {
		return errors.New("len(chatIds) must be more then 0")
	}
	err := s.repo.UpdateDaysForSubscribers(chatIds)
	return err
}

func (s *ContentService) getFavoriteAyatsInlineKeyboard(ayats []qbot.Ayat, index int) tgbotapi.InlineKeyboardMarkup {
	var keyboard tgbotapi.InlineKeyboardMarkup
	textForFavorButton, dataForFavorButtonTemplate := getTextAndDataForFavoriteButton(true)
	addToFavoriteButton := tgbotapi.NewInlineKeyboardButtonData(
		textForFavorButton,
		fmt.Sprintf(dataForFavorButtonTemplate, ayats[index].Id),
	)
	if index == 0 {
		if len(ayats) == 1 {
			keyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					addToFavoriteButton,
				),
			)
		} else {
			keyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					addToFavoriteButton,
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						fmt.Sprintf("%s >", ayats[index].GetSuraAyatNum()),
						fmt.Sprintf("getFavoriteAyat(%d)", ayats[index].Id),
					),
				),
			)
		}
	} else if index == len(ayats)-1 {
		index = index - 1
		if len(ayats) == 1 {
			keyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					addToFavoriteButton,
				),
			)
		} else {
			keyboard = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					addToFavoriteButton,
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						fmt.Sprintf("< %s", ayats[index].GetSuraAyatNum()),
						fmt.Sprintf("getFavoriteAyat(%d)", ayats[index].Id),
					),
				),
			)
		}
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				addToFavoriteButton,
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("< %s", ayats[index-1].GetSuraAyatNum()),
					fmt.Sprintf("getFavoriteAyat(%d)", ayats[index-1].Id),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s >", ayats[index+1].GetSuraAyatNum()),
					fmt.Sprintf("getFavoriteAyat(%d)", ayats[index+1].Id),
				),
			),
		)
	}
	return keyboard
}

func renderAyat(ayat qbot.Ayat) string {
	contentTemplate := "%s) %s\n\nСсылка на [источник](https://umma.ru%s)"
	return fmt.Sprintf(contentTemplate, ayat.GetSuraAyatNum(), ayat.Content, ayat.SuraLink)
}

func checkAyatInResult(query string, ayat qbot.Ayat) bool {
	switch {
	case strings.Contains(ayat.Ayat, "-"):
		return serviceNonIntAyatId(ayat.Ayat, query, "-")
	case strings.Contains(ayat.Ayat, ","):
		return serviceNonIntAyatId(ayat.Ayat, query, ",")
	case ayat.Ayat == query:
		return true
	default:
		return false
	}
}

func serviceNonIntAyatId(ayatId string, queryAyat string, separator string) bool {
	queryAyatInt, err := strconv.Atoi(queryAyat)
	if err != nil {
		return false
	}
	separatedAyatId := strings.Split(ayatId, separator)
	leftBorder, err := strconv.Atoi(strings.TrimSpace(separatedAyatId[0]))
	if err != nil {
		return false
	}
	rightBorder, err := strconv.Atoi(strings.TrimSpace(separatedAyatId[1]))
	if err != nil {
		return false
	}
	for leftBorder <= rightBorder {
		if queryAyatInt == leftBorder {
			return true
		}
		leftBorder++
	}
	return false
}
