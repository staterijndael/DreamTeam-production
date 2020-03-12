package user

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//установка аватара пользователя.
//кодировка - base64. ответственность за кодировку лежит на пользователе.
//jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//при отправке уведомления:
//uid заменяется полной информацией о пользователе.
//avatar заменяется с fileInput на fileMetaInfo, содержащий мета информацию о файле, но не контент файла
//zenrpc:avatar filtInput. контент файла в кодировке base64.
//zenrpc:return при удачном выполнении запроса возвращает мета информацию о файле.
func (s *Service) SetAvatar(
	ctx context.Context,
	avatar *common.FileInput,
) (*views.FileMetaInfo, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	dbFile, err := common.CreateDBFile(s.db, s.conf, []byte(avatar.Content))
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	me.AvatarID = dbFile.ID
	me.Avatar = *dbFile
	if err := s.db.
		Model(&models.User{}).Where("id = ?", me.ID).
		Update("avatar", dbFile.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return views.FileMetaInfoViewFromModel(dbFile), nil
}
