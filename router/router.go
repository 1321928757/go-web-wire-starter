package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/handler"
	"go-web-wire-starter/internal/mildware"
	"go-web-wire-starter/util/path"
	"path/filepath"
)

// RouterSet 路由器、处理器、服务、数据访问对象的集合
var ProviderSet = wire.NewSet(NewRouter)

func NewRouter(
	conf *config.Configuration,
	userHandler *handler.UserHandler,
	mediaHandler *handler.MediaHandler,
	corsM *mildware.Cors,
	jwtM *mildware.JWTAuth,
	recoveryM *mildware.Recovery,
) *gin.Engine {
	if conf.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	// 基本中间件
	router.Use(gin.Logger(), recoveryM.Handler())

	// 跨域处理
	router.Use(corsM.Handler())

	// 限流处理
	//router.Use(limiterM.Handler())

	rootDir := path.RootPath()
	// 前端项目静态资源
	router.StaticFile("/", filepath.Join(rootDir, "static/dist/index.html"))
	router.Static("/assets", filepath.Join(rootDir, "static/dist/assets"))
	router.StaticFile("/favicon.ico", filepath.Join(rootDir, "static/dist/favicon.ico"))
	// 其他静态资源
	router.Static("/public", filepath.Join(rootDir, "static"))
	router.Static("/storage", filepath.Join(rootDir, "storage/internal/public"))

	// 注册 api 分组路由
	groupRouter := router.Group("/api")
	setUserGroupRoutes(groupRouter, userHandler, jwtM)
	setMediaGroupRoutes(groupRouter, mediaHandler, jwtM)

	return router
}
