package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-web-wire-starter/internal/domain"
	"go-web-wire-starter/internal/pkg/request"
	"go-web-wire-starter/internal/pkg/response"
	"go-web-wire-starter/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *zap.Logger
	userService *service.UserService
	jwtService  *service.JwtService
}

func NewUserHandler(logger *zap.Logger, userService *service.UserService, jwtService *service.JwtService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
		jwtService:  jwtService,
	}
}

func (u *UserHandler) Register(c *gin.Context) {
	var form request.Register
	if err := c.ShouldBindJSON(&form); err != nil {
		response.FailByErr(c, request.GetError(form, err))
		return
	}
	user, err := u.userService.Register(c, &form)
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	tokenData, _, err := u.jwtService.CreateToken(domain.AppGuardName, user)
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	response.Success(c, tokenData)
}

func (h *UserHandler) Login(c *gin.Context) {
	var form request.Login
	if err := c.ShouldBindJSON(&form); err != nil {
		response.FailByErr(c, request.GetError(form, err))
		return
	}

	user, err := h.userService.Login(c, form.Mobile, form.Password)
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	tokenData, _, err := h.jwtService.CreateToken(domain.AppGuardName, user)
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	response.Success(c, tokenData)
}

func (h *UserHandler) Info(c *gin.Context) {
	user, err := h.userService.GetUserInfo(c, c.Keys["id"].(string))
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Logout(c *gin.Context) {
	err := h.jwtService.JoinBlackList(c, c.Keys["token"].(*jwt.Token))
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	response.Success(c, nil)
}
