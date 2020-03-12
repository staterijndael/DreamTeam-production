package services

import (
	"dt/rpc/services/admin"
	"dt/rpc/services/auth"
	"dt/rpc/services/bug"
	"dt/rpc/services/plebs"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	auth.ProviderSet,
	plebs.ProviderSet,
	admin.ProviderSet,
	bug.ProviderSet,
)
