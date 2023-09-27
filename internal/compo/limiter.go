package compo

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type LimiterManager struct {
	limiters *sync.Map
	once     *sync.Once
}

func NewLimiterManager() *LimiterManager {
	return &LimiterManager{
		limiters: &sync.Map{},
		once:     &sync.Once{},
	}
}

// Limiter 限流器
type Limiter struct {
	L           *rate.Limiter
	lastGetTime time.Time
}

// GetLimiter 获取限流器
// Param r 限流器的速率限制
// Param b 限流器的并发数量限制
// Param key 限流的资源标识 key
func (lm *LimiterManager) GetLimiter(r rate.Limit, b int, key string) *Limiter {
	// 启动一个 goroutine 定时清理过期的限流器
	lm.once.Do(func() {
		go lm.clearLimiter()
	})

	// 从 sync.Map 中获取限流器，如果不存在则新建一个
	limiter, ok := lm.limiters.Load(key)
	if ok {
		return limiter.(*Limiter)
	}

	l := &Limiter{
		L:           rate.NewLimiter(r, b),
		lastGetTime: time.Now(),
	}

	lm.limiters.Store(key, l)

	return l
}

// clearLimiter 定时清理过期的限流器
func (lm *LimiterManager) clearLimiter() {
	for {
		// 每隔一分钟清理一次过期的限流器
		time.Sleep(1 * time.Minute)

		lm.limiters.Range(func(key, value interface{}) bool {
			if time.Now().Unix()-value.(*Limiter).lastGetTime.Unix() > 180 {
				lm.limiters.Delete(key)
			}
			return true
		})
	}
}
