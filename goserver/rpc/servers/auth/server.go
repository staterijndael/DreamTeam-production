package auth

import (
	"dt/rpc/services/auth/code"
	"dt/rpc/services/auth/token"
	"dt/rpc/services/auth/user"
	"github.com/semrush/zenrpc"
)

type Server struct {
	*zenrpc.Server
}

func New(
	ts *token.Service,
	us *user.Service,
	cs *code.Service,
) *Server {
	rpc := zenrpc.NewServer(zenrpc.Options{
		ExposeSMD: true,
		AllowCORS: true,
	})

	rpc.Register("token", ts)
	rpc.Register("code", cs)
	rpc.Register("user", us)

	return &Server{
		Server: &rpc,
	}
}
