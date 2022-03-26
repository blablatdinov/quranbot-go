package telegramsdk

type UpdatesResponse struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

type Update struct {
	UpdateId int64             `json:"update_id"`
	Message  messageResultJson `json:"message"`
}
