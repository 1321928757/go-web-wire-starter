package router

import (
	"github.com/gin-gonic/gin"
	"go-web-wire-starter/internal/domain"
	"go-web-wire-starter/internal/handler"
	"go-web-wire-starter/internal/mildware"
)

// 设置用户路由
func setUserGroupRoutes(
	router *gin.RouterGroup,
	userHandler *handler.UserHandler,
	jwtM *mildware.JWTAuth,
) *gin.RouterGroup {
	group := router.Group("/user")
	group.POST("/register", userHandler.Register)
	group.POST("/login", userHandler.Login)

	authGroup := group.Group("/auth").Use(jwtM.LoginAutoHandler(domain.AppGuardName))
	{
		authGroup.GET("/info", userHandler.Info)
		authGroup.POST("/logout", userHandler.Logout)
	}

	return group
}

// 设置媒体路由
func setMediaGroupRoutes(
	router *gin.RouterGroup,
	mediaHandler *handler.MediaHandler,
	jwtM *mildware.JWTAuth,
) *gin.RouterGroup {
	group := router.Group("/media")
	authGroup := group.Group("/auth").Use(jwtM.LoginAutoHandler(domain.AppGuardName))
	{
		authGroup.POST("/image/upload", mediaHandler.ImageUpload)
		authGroup.GET("/url", mediaHandler.GetUrlById)
	}
	return group
}

// 设置验证码路由
func setCaptchaGroupRoutes(
	router *gin.RouterGroup,
	captchaHandler *handler.CaptchaHandler,
	jwtM *mildware.JWTAuth,
) *gin.RouterGroup {
	group := router.Group("/captcha")
	group.GET("/send_email", captchaHandler.SendEmailCaptcha)
	return group
}
