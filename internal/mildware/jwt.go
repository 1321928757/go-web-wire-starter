package mildware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/domain"
	cErr "go-web-wire-starter/internal/pkg/error"
	"go-web-wire-starter/internal/pkg/response"
	"go-web-wire-starter/internal/service"
	"strconv"
	"time"
)

type JWTAuth struct {
	conf *config.Configuration
	jwtS *service.JwtService
}

func NewJWTAuthM(
	conf *config.Configuration,
	jwtS *service.JwtService,
) *JWTAuth {
	return &JWTAuth{
		conf: conf,
		jwtS: jwtS,
	}
}

// LoginAutoHandler 登录校验中间件
func (m *JWTAuth) LoginAutoHandler(guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.FailByErr(c, cErr.Unauthorized("missing Authorization header"))
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &domain.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.conf.Jwt.Secret), nil
		})
		if err != nil || m.jwtS.IsInBlacklist(c, tokenStr) {
			response.FailByErr(c, cErr.Unauthorized("登录授权已失效"))
			return
		}

		claims := token.Claims.(*domain.CustomClaims)
		if claims.Issuer != guardName {
			response.FailByErr(c, cErr.Unauthorized("登录授权已失效"))
			return
		}
		// token 续签
		if int64(claims.ExpiresAt.Sub(time.Now()).Seconds()) < m.conf.Jwt.RefreshGracePeriod {
			tokenData, err := m.jwtS.RefreshToken(c, guardName, token)
			if err == nil {
				c.Header("new-token", tokenData.AccessToken)
				c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
			}
		}

		c.Set("token", token)
		c.Set("id", claims.ID)
	}
}
