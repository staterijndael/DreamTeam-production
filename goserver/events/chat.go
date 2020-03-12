package events

import "dt/models"

type ChatSentMsg struct {
	EventBase `json:"-"`
	MsgID     uint `json:"messageID"`
	MsgModel  *models.Message
}