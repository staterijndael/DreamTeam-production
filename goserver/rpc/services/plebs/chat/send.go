package chat

import (
	"bytes"
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
	"net/http"
)

//zenrpc:72 чат не найден
//zenrpc:73 вы не имеете доступа к этому чату
//zenrpc:return отправленное сообщение
func (s *Service) SendMessage(ctx context.Context, cid uint, text string) (*views.Msg, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var chat models.Chat
	if err := s.db.First(&chat, cid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.ChatNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, nil, nil)
	}

	if !chat.Community.Contains(me.ID) {
		return nil, errors.New(errors.CantAccessChat, nil, nil)
	}

	msg := models.Message{
		Text:     text,
		SenderID: me.ID,
		ChatID:   cid,
	}

	if err := s.db.Create(&msg).Error; err != nil {
		return nil, errors.New(errors.Internal, nil, nil)
	}

	s.emitter.Emit(&events.ChatSentMsg{
		EventBase: events.EventBase{Context: ctx},
		MsgID:     msg.ID,
	})

	if err := s.db.First(&msg, msg.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, nil, nil)
	}

	data := bytes.NewReader([]byte(`{
	"message":"assadasd",
	"chat":"2",
	"title":"dada"
}`))

	http.Post("localhost:8080/sendmessage", "application/json", data)

	return views.MsgFromModel(&msg), nil
}
