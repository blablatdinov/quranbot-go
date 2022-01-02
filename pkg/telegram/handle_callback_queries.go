package telegram

import (
	"errors"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleQuery(callbackQuery *tgbotapi.CallbackQuery) error {
	patterns := map[string]interface{}{
		`getFavoriteAyat\(\d+\)`:         b.swipeToFavoriteAyat,
		`getAyat\(\d+\)`:                 b.swipeToAyat,
		`addToFavorite\(\d+\)`:           b.addToFavorite,
		`removeFromFavorite\(\d+\)`:      b.removeFromFavorite,
		`setPrayerStatusToUnread\(\d+\)`: b.setPrayerStatusToUnread,
		`setPrayerStatusToRead\(\d+\)`:   b.setPrayerStatusToRead,
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

func (b *Bot) setPrayerStatusToRead(callbackQuery *tgbotapi.CallbackQuery) error {
	prayerAtUserId, err := strconv.Atoi(callbackQuery.Data[22 : len(callbackQuery.Data)-1])
	if err != nil {
		return err
	}
	keyboard, err := b.service.ChangePrayerStatus(prayerAtUserId, true)
	if err != nil {
		return err
	}
	edit := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard)
	b.bot.Send(edit)
	return nil
}

func (b *Bot) setPrayerStatusToUnread(callbackQuery *tgbotapi.CallbackQuery) error {
	prayerAtUserId, err := strconv.Atoi(callbackQuery.Data[24 : len(callbackQuery.Data)-1])
	if err != nil {
		return err
	}
	keyboard, err := b.service.ChangePrayerStatus(prayerAtUserId, false)
	if err != nil {
		return err
	}
	edit := tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard)
	b.bot.Send(edit)
	return nil
}

func (b *Bot) swipeToFavoriteAyat(callbackQuery *tgbotapi.CallbackQuery) error {
	ayatId, err := strconv.Atoi(callbackQuery.Data[16 : len(callbackQuery.Data)-1])
	if err != nil {
		return err
	}
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID
	answer, keyboard, err := b.service.GetFavoriteAyatsFromKeyboard(callbackQuery.Message.Chat.ID, ayatId)
	if err != nil {
		return err
	}
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
	state, err := b.service.GetSubscriberState(chatId)
	if err != nil {
		return err
	}
	res, keyboard, err := b.service.RemoveFromFavorite(chatId, ayatId, state)
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
	state, err := b.service.GetSubscriberState(chatId)
	if err != nil {
		return err
	}
	res, keyboard, err := b.service.AddToFavorite(chatId, ayatId, state)
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
	state, err := b.service.GetSubscriberState(chatId)
	if err != nil {
		return err
	}
	answer, keyboard, err := b.service.GetAyatById(chatId, ayatId, state)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewEditMessageText(chatId, messageId, answer)
	msg.ParseMode = "markdown"
	b.bot.Send(msg)
	edit := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)
	b.bot.Send(edit)
	return nil
}
