package admin

import (
	"dt/rpc/services/admin/fns"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	fns.New,
)
