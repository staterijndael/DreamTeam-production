package bug

import (
	"dt/rpc/services/bug/report"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	report.New,
)
