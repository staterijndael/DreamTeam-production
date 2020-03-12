package chat

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
	"time"
)

//возращает до 50 сообщений чата, по умолчанию возращаются последние 50, при указании времени - до этого времени
//zenrpc:after=0 время до которого нужны сообщения
//zenrpc:72 чат с таким id не найден
//zenrpc:73 у вас нет доступа к этому чату
//zenrpc: return сообщение из чата
func (s *Service) LoadPageOfMessages(
	ctx context.Context,
	cid uint,
	after int64,
) (messages []*views.Msg, err *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var chat models.Chat
	if err := s.db.First(&chat, cid).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.Internal, nil, nil)
		}

		return nil, errors.New(errors.ChatNotFound, err, nil)
	}

	if !chat.Community.Contains(me.ID) {
		return nil, errors.New(errors.CantAccessChat, nil, nil)
	}

	var msgs []*models.Message
	afterTime := time.Unix(after, 0)
	var scope func(db *gorm.DB) *gorm.DB
	if after == 0 {
		scope = scopes.GetMessagesByChat(cid)
	} else {
		scope = scopes.GetMessagesByChatBeforeTime(cid, &afterTime)
	}

	if err := s.db.Scopes(scope).Find(&msgs).Error; err != nil {
		return nil, errors.New(errors.Internal, nil, nil)
	}

	for _, msg := range msgs {
		messages = append(messages, views.MsgFromModel(msg))
	}

	return messages, nil
}
