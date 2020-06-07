package app

import "encoding/json"

type onSendMsg struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Email       string `json:"email"`
}

func onSendMarshaling(event *onSendMsg) []byte {
	marshaled, _ := json.Marshal(event)
	return marshaled
}
