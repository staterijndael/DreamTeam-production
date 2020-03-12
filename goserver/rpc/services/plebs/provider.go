package plebs

import (
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
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	rating.New,
	user.New,
	notification.New,
	group_join.New,
	group_creation.New,
	group.New,
	org.New,
	org_join.New,
	fns.New,
	chat.New,
	wall.New,
)
