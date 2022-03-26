package telegramsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

func (b *Bot) SendMessage(chatId int64, text string) (Message, error) {
	response, err := http.Get(b.getUrl(fmt.Sprintf("/sendMessage?chat_id=%d&text=%s", chatId, text)))
	if err != nil {
		return Message{}, err
	}
	defer response.Body.Close()
	var messageJson MessageJson
	if err = json.NewDecoder(response.Body).Decode(&messageJson); err != nil {
		return Message{}, err
	}
	message := Message{
		Ok: messageJson.Ok,
		Result: messageResult{
			MessageId: messageJson.Result.MessageId,
			From:      messageJson.Result.From,
			Chat:      messageJson.Result.Chat,
			Date:      time.Unix(messageJson.Result.Date, 0),
			Text:      messageJson.Result.Text,
		},
	}
	return message, nil
}

func (b *Bot) getUrl(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s%s", b.Token, method)
}
