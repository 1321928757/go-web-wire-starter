package email

import (
	"github.com/jordan-wright/email"
	"go-web-wire-starter/config"
	"go.uber.org/zap"
	"net/smtp"
)

func NewEmailPool(config *config.Configuration, logger *zap.Logger) *email.Pool {
	// 初始化邮箱连接池
	emailConfig := config.Email
	pool, err := email.NewPool(
		emailConfig.Host+":"+emailConfig.Port,
		emailConfig.MaxConnection,
		smtp.PlainAuth("", emailConfig.SenderEmail, emailConfig.SenderPassword,
			emailConfig.Host),
	)
	if err != nil {
		logger.Fatal("初始化邮件连接池失败, 错误:", zap.Any("err", err))
	}
	return pool
}
