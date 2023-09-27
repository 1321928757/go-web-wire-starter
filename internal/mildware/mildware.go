package mildware

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCorsM, NewRecoveryM, NewJWTAuthM, NewLimiterM)
