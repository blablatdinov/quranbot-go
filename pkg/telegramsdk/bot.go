package telegramsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const updatesTimeout = 1 * time.Second

type Bot struct {
	Token string
}

func NewBot(token string) *Bot {
	return &Bot{token}
}

func (b *Bot) GetMe() (GetMeStruct, error) {
	response, err := http.Get(b.getUrl("/getMe"))
	if err != nil {
		return GetMeStruct{}, err
	}
	defer response.Body.Close()
	var botData GetMeStruct
	if err = json.NewDecoder(response.Body).Decode(&botData); err != nil {
		return GetMeStruct{}, err
	}
	return botData, nil
}

func (b *Bot) sendMessage(chatId int64, text string, keyboard string) (Message, error) {
	url := b.getUrl(fmt.Sprintf("/sendMessage?chat_id=%d&text=%s&reply_markup=%s", chatId, text, keyboard))
	response, err := http.Get(url)
	if err != nil {
		return Message{}, err
	}
	defer response.Body.Close()
	var messageJson MessageJson
	if err = json.NewDecoder(response.Body).Decode(&messageJson); err != nil {
		return Message{}, err
	}
	return messageJsonToMessage(messageJson), nil
}

func (b *Bot) SendMessage(chatId int64, text string) (Message, error) {
	defaultKeyboard := getDefaultKeyboardJson()
	return b.sendMessage(chatId, text, defaultKeyboard)
}

func (b *Bot) SendMessageWithKeyboard(chatId int64, text string, keyboard string) (Message, error) {
	return b.sendMessage(chatId, text, keyboard)
}

func getDefaultKeyboardJson() string {
	defaultKeyboard := ReplyKeyboardMarkup{
		Keyboard: [][]ReplyKeyboardButton{
			{
				{"🎧 Подкасты"},
			},
			{
				{"🕋 Время намаза"},
			},
			{
				{"🌟 Избранное"}, {"🔍 Найти аят"},
			},
		},
	}
	keyboardJson, err := json.Marshal(defaultKeyboard)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(keyboardJson)
}

func logMessage(message Message) {
	if len(message.Text) > 50 {
		message.Text = message.Text[:50]
	}
	log.Printf("Getting message id:%d, chat_id: %d, text: %s...", message.MessageId, message.Chat.Id, message.Text)
}

func (b *Bot) GetUpdatesChan() chan Message {
	log.Printf("Bot started on updates. Timeout before getUpdates: %d\n", updatesTimeout)
	updatesChan := make(chan Message)
	go func(updatesChan chan Message) {
		var offset int64 = 0
		for {
			messages, lastUpdateId, _ := b.GetUpdates(offset + 1)
			offset = lastUpdateId
			for _, message := range messages {
				logMessage(message)
				updatesChan <- message
			}
			time.Sleep(updatesTimeout)
		}
	}(updatesChan)
	return updatesChan
}

func (b *Bot) DeleteWebhook() error {
	response, err := http.Get(b.getUrl("/deleteWebhook"))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Can't delete webhook: %s", string(bytes)))
	}
	return nil
}

func (b *Bot) SetWebhook(host string) error {
	url := b.getUrl(fmt.Sprintf("/setWebhook?url=%s", host))
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("Can't set webhook: %s", string(bytes)))
	}
	return nil
}

func (b *Bot) RunWebhookServer(host string) (chan Message, error) {
	updatesChan := make(chan Message)
	if err := b.DeleteWebhook(); err != nil {
		return updatesChan, err
	}
	if err := b.SetWebhook(host); err != nil {
		return updatesChan, err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var update Update
		updatesChan <- messageJsonResultToMessage(update.Message)
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Fatal(err.Error())
		}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(bytes))
	})
	go func() {
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatal(err.Error())
		}
	}()
	log.Println("Bot started on webhook")
	return updatesChan, nil
}

func (b *Bot) GetUpdates(offset int64) ([]Message, int64, error) {
	url := b.getUrl(fmt.Sprintf("/getUpdates?offset=%d", offset))
	response, err := http.Get(url)
	if err != nil {
		return []Message{}, 0, err
	}
	defer response.Body.Close()
	var updates UpdatesResponse
	if err = json.NewDecoder(response.Body).Decode(&updates); err != nil {
		return []Message{}, 0, err
	}
	var messages []Message
	for _, update := range updates.Updates {
		messages = append(messages, messageJsonResultToMessage(update.Message))
	}
	var lastUpdateIndex int64 = 0
	if len(updates.Updates) > 0 {
		lastUpdateIndex = updates.Updates[(len(updates.Updates) - 1)].UpdateId
	}
	return messages, lastUpdateIndex, nil
}

func (b *Bot) SendAnswer(answer Answer) (Message, error) {
	return Message{}, nil
}

func (b *Bot) getUrl(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s%s", b.Token, method)
}
