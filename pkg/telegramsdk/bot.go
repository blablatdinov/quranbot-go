package telegramsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (b *Bot) getUrl(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s%s", b.Token, method)
}
