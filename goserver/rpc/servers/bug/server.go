package bug

import (
	"dt/rpc/servers/common"
	"dt/rpc/services/bug/report"
	"github.com/semrush/zenrpc"
)

type Server struct {
	*zenrpc.Server
}

func New(
	rs *report.Service,
) *Server {
	rpc := zenrpc.NewServer(zenrpc.Options{
		ExposeSMD: true,
		AllowCORS: true,
	})

	rpc.Register("report", rs)
	rpc.Use(common.Logger)

	return &Server{
		Server: &rpc,
	}
}
