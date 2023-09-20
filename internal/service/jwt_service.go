package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/compo"
	"go-web-wire-starter/internal/dao"
	"go-web-wire-starter/internal/domain"
	cErr "go-web-wire-starter/internal/pkg/error"
	"go.uber.org/zap"
	"time"
)

type JwtService struct {
	logger      *zap.Logger
	conf        *config.Configuration
	userService *UserService
	jwtDao      *dao.JwtDao
	lock        *compo.LockBuilder
}

func NewJwtService(logger *zap.Logger, conf *config.Configuration, userService *UserService,
	jwtDao *dao.JwtDao, lock *compo.LockBuilder) *JwtService {
	return &JwtService{
		logger:      logger,
		conf:        conf,
		userService: userService,
		jwtDao:      jwtDao,
		lock:        lock,
	}
}

type JwtRepo interface {
	JoinBlackList(ctx context.Context, tokenStr string, joinUnix int64, expires time.Duration) error
	GetBlackJoinUnix(ctx context.Context, tokenStr string) (int64, error)
}

func (s *JwtService) CreateToken(GuardName string, user domain.JwtUser) (*domain.TokenOutPut, *jwt.Token, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		domain.CustomClaims{
			Key: GuardName,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    GuardName,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.conf.Jwt.JwtTtl))),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * -1000)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        user.GetUid(),
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(s.conf.Jwt.Secret))
	if err != nil {
		return nil, nil, cErr.BadRequest("create token error:" + err.Error())
	}

	return &domain.TokenOutPut{
		AccessToken: tokenStr,
		ExpiresIn:   int(s.conf.Jwt.JwtTtl),
	}, token, nil
}

func (s *JwtService) JoinBlackList(ctx *gin.Context, token *jwt.Token) error {
	nowUnix := time.Now().Unix()
	timer := token.Claims.(*domain.CustomClaims).ExpiresAt.Sub(time.Now())
	fmt.Println("JoinBlackList timer", timer)

	if err := s.jwtDao.JoinBlackList(ctx, token.Raw, nowUnix, timer); err != nil {
		s.logger.Error(err.Error())
		return cErr.BadRequest("登出失败")
	}

	return nil
}

func (s *JwtService) IsInBlacklist(ctx *gin.Context, tokenStr string) bool {
	joinUnix, err := s.jwtDao.GetBlackJoinUnix(ctx, tokenStr)
	if err != nil {
		return false
	}

	if time.Now().Unix()-joinUnix < s.conf.Jwt.JwtBlacklistGracePeriod {
		return false
	}
	return true
}

func (s *JwtService) GetUserInfo(ctx *gin.Context, guardName, id string) (domain.JwtUser, error) {
	switch guardName {
	case domain.AppGuardName:
		return s.userService.GetUserInfo(ctx, id)
	default:
		return nil, cErr.BadRequest("guard " + guardName + " does not exist")
	}
}

func (s *JwtService) RefreshToken(ctx *gin.Context, guardName string, token *jwt.Token) (*domain.TokenOutPut, error) {
	idStr := token.Claims.(*domain.CustomClaims).ID

	lock := s.lock.NewLock(ctx, "refresh_token_lock:"+idStr, s.conf.Jwt.JwtBlacklistGracePeriod)
	if lock.Get() {
		user, err := s.GetUserInfo(ctx, guardName, idStr)
		if err != nil {
			s.logger.Error(err.Error())
			lock.Release()
			return nil, err
		}

		tokenData, _, err := s.CreateToken(guardName, user)
		if err != nil {
			lock.Release()
			return nil, err
		}

		err = s.JoinBlackList(ctx, token)
		if err != nil {
			lock.Release()
			return nil, err
		}

		return tokenData, nil
	}

	return nil, cErr.BadRequest("系统繁忙")
}
