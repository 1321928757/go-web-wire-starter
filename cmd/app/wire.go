//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/compo"
	"go-web-wire-starter/internal/dao"
	"go-web-wire-starter/internal/handler"
	"go-web-wire-starter/internal/mildware"
	"go-web-wire-starter/internal/service"
	"go-web-wire-starter/router"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

// wireApp init application.
func wireApp(*config.Configuration, *lumberjack.Logger, *zap.Logger) (*App, func(), error) {
	panic(
		wire.Build(
			compo.ProviderSet,
			mildware.ProviderSet,
			dao.ProviderSet,
			service.ProviderSet,
			handler.ProviderSet,
			router.ProviderSet,
			newHttpServer,
			newApp,
		),
	)
}
