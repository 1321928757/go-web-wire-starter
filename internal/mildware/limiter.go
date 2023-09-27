package mildware

import (
	"github.com/gin-gonic/gin"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/compo"
	cErr "go-web-wire-starter/internal/pkg/error"
	"go-web-wire-starter/internal/pkg/response"
	"golang.org/x/time/rate"
	"time"
)

type Limiter struct {
	lm     *compo.LimiterManager
	config *config.Configuration
}

func NewLimiterM(lm *compo.LimiterManager, config *config.Configuration) *Limiter {
	return &Limiter{
		lm:     lm,
		config: config,
	}
}

func (m *Limiter) Handler(key ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var limiterKey string
		//如果传入的key不为空，则使用传入的key，如果传入的key为空，则使用用户token作为限流key
		if len(key) > 0 && len(key[0]) > 0 {
			limiterKey = key[0]
		} else {
			// 如果token为空,则使用客户端IP作为key
			limiterKey = ctx.GetString("token")
			if len(limiterKey) == 0 {
				limiterKey = ctx.ClientIP()
			}
		}

		// 获取限流器对象， 参数rate.Every(50*time.Millisecond)表示每50毫秒向令牌桶中放入一个令牌
		// 300为令牌桶的容量，limiterKey为限流器的key
		l := m.lm.GetLimiter(rate.Every(time.Duration(m.config.Limiter.Rate)*time.Millisecond),
			m.config.Limiter.Capacity, limiterKey)

		if !l.L.Allow() {
			response.FailByErr(ctx, cErr.TooManyRequestsErr("您的访问过于频繁，请稍候重试"))
			return
		}
	}
}
