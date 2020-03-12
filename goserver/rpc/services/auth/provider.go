package auth

import (
	"dt/rpc/services/auth/code"
	"dt/rpc/services/auth/token"
	"dt/rpc/services/auth/user"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	code.New,
	token.New,
	user.New,
)
