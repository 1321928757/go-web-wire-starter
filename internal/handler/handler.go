package handler

import (
	"github.com/google/wire"
)

// ProviderSet Provider对象集合
var ProviderSet = wire.NewSet(NewUserHandler, NewMediaHandler, NewCaptchaHandler)
