package admin

import (
	"dt/rpc/servers/common"
	"dt/rpc/services/admin/fns"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

type Server struct {
	*zenrpc.Server
	sqlStore *gorm.DB
}

func New(
	db *gorm.DB,
	fns *fns.Service,
) *Server {
	rpc := zenrpc.NewServer(zenrpc.Options{
		ExposeSMD: true,
		AllowCORS: true,
		Upgrader:  common.Upgrader,
	})

	rpc.Register("fns", fns)
	rpc.Use(
		common.Logger,
		common.RequestBuilder,
	)

	return &Server{
		Server:   &rpc,
		sqlStore: db,
	}
}
