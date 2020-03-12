package managers

import (
	"dt/managers/connections"
	"dt/managers/eventEmitter"
	"dt/managers/fns"
	"dt/managers/rating"
	"dt/managers/sms"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	connections.NewManager,
	sms.NewManager,
	connections.NewDashConnectionManager,
	fns.NewManager,
	rating.New,
	eventEmitter.NewEventEmitter,
)
