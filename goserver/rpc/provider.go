package rpc

import (
	"dt/rpc/servers"
	"dt/rpc/services"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	services.ProviderSet,
	servers.ProviderSet,
)
