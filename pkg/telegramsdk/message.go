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

type Message struct {
	Ok     bool
	Result messageResult
}

type messageResult struct {
	MessageId int
	From      from
	Chat      chat
	Date      time.Time
	Text      string
}
