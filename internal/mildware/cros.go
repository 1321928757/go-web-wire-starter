package mildware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Cors struct {
}

func NewCorsM() *Cors {
	return &Cors{}
}

// Handler 跨域中间件
func (m *Cors) Handler() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"New-Token", "New-Expires-In", "Content-Disposition"}

	return cors.New(config)
}
