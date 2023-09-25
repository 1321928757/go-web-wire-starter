package compo

import (
	"github.com/google/wire"
	"go-web-wire-starter/internal/compo/email"
	"go-web-wire-starter/internal/compo/storage"
	"go-web-wire-starter/internal/compo/storage/cos"
	"go-web-wire-starter/internal/compo/storage/kodo"
	"go-web-wire-starter/internal/compo/storage/local"
	"go-web-wire-starter/internal/compo/storage/oss"
)

// ProviderSet is compo providers.
var ProviderSet = wire.NewSet(
	NewSonyFlake,
	NewLockBuilder,
	cos.NewCosDriver,
	local.NewLocalDriver,
	oss.NewOssDriver,
	kodo.NewKodoDriver,
	storage.NewStorage,
	email.NewEmailPool,
	email.NewEmailDriver,
	NewCaptchaCompo,
)
