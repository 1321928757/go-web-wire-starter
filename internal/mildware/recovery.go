package mildware

import (
	"github.com/gin-gonic/gin"
	"go-web-wire-starter/internal/pkg/response"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Recovery struct {
	loggerWriter *lumberjack.Logger
}

func NewRecoveryM(loggerWriter *lumberjack.Logger) *Recovery {
	return &Recovery{
		loggerWriter: loggerWriter,
	}
}

// Handler Recovery 中间件
func (m *Recovery) Handler() gin.HandlerFunc {
	return gin.RecoveryWithWriter(
		m.loggerWriter,
		response.ServerError,
	)
}
