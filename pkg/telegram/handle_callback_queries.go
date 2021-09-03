package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func (b *Bot) handleQuery(callbackQuery *tgbotapi.CallbackQuery) error {
	patterns := map[string]interface{}{
		`getFavoriteAyat\(\d+\)`:    b.swipeToFavoriteAyat,
		`getAyat\(\d+\)`:            b.swipeToAyat,
		`addToFavorite\(\d+\)`:      b.addToFavorite,
		`removeFromFavorite\(\d+\)`: b.removeFromFavorite,
	}
	for pattern, handler := range patterns {
		if path(pattern, callbackQuery.Data) {
			err := handler.(func(callbackQuery *tgbotapi.CallbackQuery) error)(callbackQuery)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("unknow pattern")
}

func (b *Bot) swipeToFavoriteAyat(callbackQuery *tgbotapi.CallbackQuery) error {
	ayatId, err := strconv.Atoi(callbackQuery.Data[16 : len(callbackQuery.Data)-1])
	if err != nil {
		return err
	}
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID
	answer, keyboard, err := b.service.GetFavoriteAyatsFromKeyboard(callbackQuery.Message.Chat.ID, ayatId)
	msg := tgbotapi.NewEditMessageText(chatId, messageId, answer)
	msg.ParseMode = "markdown"
	b.bot.Send(msg)
	edit := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)
	b.bot.Send(edit)
	return nil
}

func (b *Bot) removeFromFavorite(callbackQuery *tgbotapi.CallbackQuery) error {
	chatId := callbackQuery.Message.Chat.ID
	ayatId, err := strconv.Atoi(callbackQuery.Data[19 : len(callbackQuery.Data)-1])
	if err != nil {
		return err
	}
	res, keyboard, err := b.service.RemoveFromFavorite(chatId, ayatId)
	if err != nil {
		return err
	}
	_, err = b.bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: callbackQuery.ID,
		Text:            res,
	})
	edit := tgbotapi.NewEditMessageReplyMarkup(chatId, callbackQuery.Message.MessageID, keyboard)
	b.bot.Send(edit)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) addToFavorite(callbackQuery *tgbotapi.CallbackQuery) error {
	ayatId, err := strconv.Atoi(callbackQuery.Data[14 : len(callbackQuery.Data)-1])
	if err != nil {
		return err
	}
	chatId := callbackQuery.Message.Chat.ID
	res, keyboard, err := b.service.AddToFavorite(chatId, ayatId)
	if err != nil {
		return err
	}
	_, err = b.bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: callbackQuery.ID,
		Text:            res,
	})
	edit := tgbotapi.NewEditMessageReplyMarkup(chatId, callbackQuery.Message.MessageID, keyboard)
	b.bot.Send(edit)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) swipeToAyat(callbackQuery *tgbotapi.CallbackQuery) error {
	ayatId, err := strconv.Atoi(callbackQuery.Data[8 : len(callbackQuery.Data)-1])
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID
	if err != nil {
		return err
	}
	answer, keyboard, err := b.service.GetAyatById(chatId, ayatId)
	msg := tgbotapi.NewEditMessageText(chatId, messageId, answer)
	msg.ParseMode = "markdown"
	b.bot.Send(msg)
	edit := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)
	b.bot.Send(edit)
	return nil
}
