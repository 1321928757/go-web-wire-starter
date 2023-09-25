package handler

import (
	"github.com/gin-gonic/gin"
	"go-web-wire-starter/internal/pkg/request"
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

// 获取图像点击验证码
func (h *CaptchaHandler) GetClickImgCaptcha(c *gin.Context) {
	captcha, err := h.captchaService.GetImgClickCaptcha()
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}
	response.Success(c, captcha)
}

// 校验图像点击验证码
func (h *CaptchaHandler) CheckClickImgCaptcha(c *gin.Context) {
	var form request.CaptchaClickParams
	if err := c.ShouldBindJSON(&form); err != nil {
		response.FailByParams(c, request.GetErrorMsg(form, err))
		return
	}
	token, result, err := h.captchaService.CheckImgClickCaptcha(form.Key, form.Dots)
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	if result {
		response.Success(c, token)
	} else {
		response.FailByBussiness(c, "人机校验不通过")
	}
}
