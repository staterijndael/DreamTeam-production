package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type ChatSentMsg struct {
	notificationBase
	Msg *models.Message `json:"msg"`
} //notifier

func (csm *ChatSentMsg) loadReceivers() {
	csm.receivers = append(csm.receivers, csm.Msg.Chat.Community.MembersIDs()...)
}

func (csm *ChatSentMsg) loadDashReceivers() {}
func (csm *ChatSentMsg) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.chatsentmessage",
		Data: csm.View(),
	}
}

func (csm *ChatSentMsg) View() interface{} {
	return &struct {
		ID   uint       `json:"id"`
		Msg  *views.Msg `json:"message"`
		Seen *bool      `json:"seen,omitempty"`
	}{
		ID:   csm.GetModel().ID,
		Msg:  views.MsgFromModel(csm.Msg),
		Seen: csm.seen,
	}
}

func (csm *ChatSentMsg) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.ChatSentMsg)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = csm.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	if _, err := saveUNS(db, n, csm.Msg.Chat.Community.MembersIDs()); err != nil {
		return err
	}

	return nil
}

func (csm *ChatSentMsg) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.ChatSentMsg
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return csm.LoadWithEvent(db, e, n)
}

func (csm *ChatSentMsg) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.ChatSentMsg
	var ok bool
	if event, ok = _event.(*events.ChatSentMsg); !ok {
		return WrongEventErr
	}

	var msg models.Message
	if err := db.Preload("Chat").First(&msg, event.MsgID).Error; err != nil {
		return err
	}

	csm.Model = model
	csm.Msg = &msg
	csm.loadReceivers()
	csm.loadDashReceivers()

	return nil
}
