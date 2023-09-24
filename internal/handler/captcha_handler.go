package handler

import (
	"github.com/gin-gonic/gin"
	"go-web-wire-starter/internal/pkg/response"
	"go-web-wire-starter/internal/service"
	"go.uber.org/zap"
)

type CaptchaHandler struct {
	logger         *zap.Logger
	captchaService *service.CaptchaService
}

func NewCaptchaHandler(logger *zap.Logger, captchaService *service.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{
		logger:         logger,
		captchaService: captchaService,
	}
}

// 发送邮箱验证码
func (h *CaptchaHandler) SendEmailCaptcha(c *gin.Context) {
	// 获取邮箱参数
	email := c.Query("email")
	if email == "" {
		response.FailByParams(c, "请求参数错误")
		return
	}

	// 发送验证码
	err := h.captchaService.SendEmailCaptcha(email)
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	response.Success(c, "ok")
}
