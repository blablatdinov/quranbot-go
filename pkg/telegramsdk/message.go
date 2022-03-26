package telegramsdk

import "time"

// {
// 	"ok":true,
// 	"result":{
// 	  "message_id":19605,
// 	  "from":{
// 		"id":452230948,
// 		"is_bot":true,
// 		"first_name":"WokeUpSmiled",
// 		"username":"WokeUpSmiled_bot"
// 	  },
// 	  "chat":{
// 		"id":358610865,
// 		"first_name":"\u0410\u043b\u043c\u0430\u0437",
// 		"last_name":"\u0418\u043b\u0430\u043b\u0435\u0442\u0434\u0438\u043d\u043e\u0432",
// 		"username":"ilaletdinov",
// 		"type":"private"
// 	  },
// 	  "date":1648284215,
// 	  "text":"asdf"
// 	}
//   }

type MessageJson struct {
	Ok     bool              `json:"ok"`
	Result messageResultJson `json:"result"`
}

type messageResultJson struct {
	MessageId int    `json:"message_id"`
	From      from   `json:"from"`
	Chat      chat   `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}

type from struct {
	Id        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type chat struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type MessageResponse struct {
	Ok     bool
	Result Message
}

type Message struct {
	MessageId int
	From      from
	Chat      chat
	Date      time.Time
	Text      string
}

func messageJsonToMessage(messageJson MessageJson) Message {
	return Message{
		MessageId: messageJson.Result.MessageId,
		From:      messageJson.Result.From,
		Chat:      messageJson.Result.Chat,
		Date:      time.Unix(messageJson.Result.Date, 0),
		Text:      messageJson.Result.Text,
	}
}

func messageJsonResultToMessage(messageJson messageResultJson) Message {
	return Message{
		MessageId: messageJson.MessageId,
		From:      messageJson.From,
		Chat:      messageJson.Chat,
		Date:      time.Unix(messageJson.Date, 0),
		Text:      messageJson.Text,
	}
}
