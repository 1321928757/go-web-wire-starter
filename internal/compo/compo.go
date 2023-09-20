package compo

import (
	"github.com/google/wire"
)

// ProviderSet is compo providers.
var ProviderSet = wire.NewSet(
	NewSonyFlake,
	NewLockBuilder,
)
