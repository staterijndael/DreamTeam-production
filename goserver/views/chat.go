package views

import "dt/models"

type Chat struct {
	ID        uint       `json:"id"`
	Community *Community `json:"community"`
}

func ChatFromModel(c *models.Chat) *Chat {
	return &Chat{
		ID:        c.ID,
		Community: CommunityFromModelShort(&c.Community),
	}
}

type Msg struct {
	ID     uint  `json:"id"`
	SentAt int64 `json:"sentAt"`
	ChatID uint  `json:"chatID"`
	Sender *User `json:"sender"`
	Text string `json:"text"`
}

func MsgFromModel(m *models.Message) *Msg {
	return &Msg{
		ID:     m.ID,
		Text: m.Text,
		SentAt: m.CreatedAt.Unix(),
		ChatID: m.ChatID,
		Sender: UserViewFromModel(&m.Sender),
	}
}
