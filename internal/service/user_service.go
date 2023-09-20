package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-web-wire-starter/internal/dao"
	"go-web-wire-starter/internal/domain"
	cErr "go-web-wire-starter/internal/pkg/error"
	"go-web-wire-starter/internal/pkg/request"
	"go-web-wire-starter/util/hash"
	"go.uber.org/zap"
	"strconv"
)

type UserRepo interface {
	FindByID(context.Context, uint64) (*domain.User, error)
	FindByMobile(context.Context, string) (*domain.User, error)
	Create(context.Context, *domain.User) (*domain.User, error)
}

type UserService struct {
	logger  *zap.Logger
	userDao *dao.UserDao
}

func NewUserService(logger *zap.Logger, userDao *dao.UserDao) *UserService {
	return &UserService{
		logger:  logger,
		userDao: userDao,
	}
}

// Register 注册
func (s *UserService) Register(ctx *gin.Context, param *request.Register) (*domain.User, error) {
	user, _ := s.userDao.FindByMobile(ctx, param.Mobile)
	if user != nil {
		return nil, cErr.BadRequest("手机号码已存在")
	}
	u, err := s.userDao.Create(ctx, &domain.User{
		Name:     param.Name,
		Mobile:   param.Mobile,
		Password: hash.BcryptMake(param.Password),
	})
	if err != nil {
		return nil, cErr.BadRequest("注册用户失败")
	}

	return u, nil
}

// Login 登录
func (s *UserService) Login(ctx *gin.Context, mobile, password string) (*domain.User, error) {
	u, err := s.userDao.FindByMobile(ctx, mobile)
	if err != nil || !hash.BcryptMakeCheck(password, u.Password) {
		return nil, cErr.BadRequest("用户名不存在或密码错误")
	}

	return u, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx *gin.Context, idStr string) (*domain.User, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, cErr.NotFound("数据ID错误")
	}
	u, err := s.userDao.FindByID(ctx, id)
	if err != nil {
		return nil, cErr.NotFound("数据不存在", cErr.USER_NOT_FOUND)
	}

	return u, nil
}
