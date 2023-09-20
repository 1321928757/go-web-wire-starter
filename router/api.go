package router

import (
	"github.com/gin-gonic/gin"
	"go-web-wire-starter/internal/domain"
	"go-web-wire-starter/internal/handler"
	"go-web-wire-starter/internal/mildware"
)

// 设置用户api路由
func setUserGroupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	jwtM *mildware.JWTAuth,
) *gin.RouterGroup {
	group := router.Group("/user")
	group.POST("/register", userHandler.Register)
	group.POST("/login", userHandler.Login)

	authGroup := group.Group("/auto").Use(jwtM.LoginAutoHandler(domain.AppGuardName))
	{
		authGroup.GET("/info", userHandler.Info)
		authGroup.POST("/logout", userHandler.Logout)
	}

	return group
}
