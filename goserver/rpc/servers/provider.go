package servers

import (
	"dt/rpc/servers/admin"
	"dt/rpc/servers/auth"
	"dt/rpc/servers/bug"
	"dt/rpc/servers/plebs"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	admin.New,
	auth.New,
	plebs.New,
	bug.New,
)
