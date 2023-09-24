package service

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/compo/email"
	"go-web-wire-starter/util/str"
	"go.uber.org/zap"
	"time"
)

// 验证码服务层
type CaptchaService struct {
	conf  *config.Configuration
	log   *zap.Logger
	email *email.EmailDriver
	rdb   *redis.Client
}

func NewCaptchaService(conf *config.Configuration, log *zap.Logger,
	email *email.EmailDriver, rdb *redis.Client) *CaptchaService {
	return &CaptchaService{conf: conf, log: log, email: email, rdb: rdb}
}

// 发送邮箱验证码
func (s *CaptchaService) SendEmailCaptcha(email string) error {
	// 判断验证码发送时长间隔，防止频繁发送
	ttl := s.rdb.TTL(context.Background(), s.conf.Captcha.CaptchaPrefix+":"+email)
	if ttl.Val().Seconds() > s.conf.Captcha.EmailExpire.Seconds()-s.conf.Captcha.EmailInterval.Seconds() {
		return errors.New("验证码发送过于频繁,请稍后再试")
	}

	// 生成验证码
	captcha := str.RandString(s.conf.Captcha.EmailNumber)

	// 发送邮件
	title := "测试验证码"
	content := "您本次操作的验证码为:" + captcha + ", 有效期为" + s.conf.Captcha.EmailExpire.String()
	err := s.email.SendRegisterMail(email, title, content)
	if err != nil {
		return err
	}

	// 保存验证码到redis缓存中
	err = s.rdb.Set(context.Background(), s.conf.Captcha.CaptchaPrefix+":"+email,
		captcha, s.conf.Captcha.EmailExpire*time.Second).Err()
	return err
}
