package plebs

import (
	"dt/managers/connections"
	"dt/rpc/servers/common"
	"dt/rpc/services/plebs/chat"
	"dt/rpc/services/plebs/fns"
	"dt/rpc/services/plebs/group"
	"dt/rpc/services/plebs/group_creation"
	"dt/rpc/services/plebs/group_join"
	"dt/rpc/services/plebs/notification"
	"dt/rpc/services/plebs/org"
	"dt/rpc/services/plebs/org_join"
	"dt/rpc/services/plebs/rating"
	"dt/rpc/services/plebs/user"
	"dt/rpc/services/plebs/wall"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

type Server struct {
	*zenrpc.Server
	sqlStore *gorm.DB
	ucm      *connections.UserConnectionsManager
	dcm      *connections.DashConnectionsManager
}

func New(
	db *gorm.DB,
	us *user.Service,
	os *org.Service,
	gcs *group_creation.Service,
	gs *group.Service,
	gjr *group_join.Service,
	ns *notification.Service,
	rs *rating.Service,
	ojr *org_join.Service,
	wall *wall.Service,
	fns *fns.Service,
	cs *chat.Service,
	ucm *connections.UserConnectionsManager,
	dcm *connections.DashConnectionsManager,
) *Server {
	rpc := zenrpc.NewServer(zenrpc.Options{
		ExposeSMD: true,
		AllowCORS: true,
		Upgrader:  common.Upgrader,
	})

	rpc.Register("rating", rs)
	rpc.Register("user", us)
	rpc.Register("organization", os)
	rpc.Register("groupCreation", gcs)
	rpc.Register("group", gs)
	rpc.Register("groupJoinRequest", gjr)
	rpc.Register("orgJoinRequest", ojr)
	rpc.Register("notification", ns)
	rpc.Register("fns", fns)
	rpc.Register("wall", wall)
	rpc.Register("chat", cs)
	rpc.Use(
		common.Logger,
		common.RequestBuilder,
		common.NicknameChecker,
	)

	return &Server{
		Server:   &rpc,
		sqlStore: db,
		ucm:      ucm,
		dcm:      dcm,
	}
}
