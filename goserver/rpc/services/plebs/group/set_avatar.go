package group

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//установка аватара группы.
//кодировка - base64. ответственность за кодировку лежит на пользователе.
//jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя,
// а так же слинкованным к орг-ии пользователям.
//zenrpc:gid id группы. при уведомлении сменяется на полную информацию о группе.
//zenrpc:avatar filtInput. контент файла в кодировке base64. при отправке уведомления заменяется с fileInput на fileMetaInfo, содержащий мета информацию о файле, но не контент файла
//zenrpc:23 group not found. группа с данным id не найдена.
//zenrpc:87 только директор и админы орг-ии, а также админ группы имеют право на данную операцию
//zenrpc:return при удачном выполнении запроса возвращает мета информацию о файле.
func (s *Service) SetAvatar(
	ctx context.Context,
	gid uint,
	avatar *common.FileInput,
) (*views.FileMetaInfo, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var group models.Group
	if err := s.db.First(&group, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if me.ID != group.Organization.DirectorID && !group.Organization.Admins.Contains(me.ID) && me.ID != group.AdminID {
		return nil, errors.New(errors.OnlyDirectorOrgAdminOrGroupAdminMaySetGroupAvatar, nil, nil)
	}

	dbFile, errDBFile := common.CreateDBFile(s.db, s.conf, []byte(avatar.Content))
	if errDBFile != nil {
		return nil, errors.New(errors.Internal, errDBFile, nil)
	}

	if err := s.db.Model(&models.Group{}).Where("id = ?", group.ID).Update("avatar", dbFile.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return views.FileMetaInfoViewFromModel(dbFile), nil
}
