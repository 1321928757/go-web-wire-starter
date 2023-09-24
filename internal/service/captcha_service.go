package service

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/compo/email"
	"go-web-wire-starter/util/str"
	"go.uber.org/zap"
	"strconv"
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
	// 如果验证码存在，且有效期大于配置的有效期减去间隔时间，则提示发送过于频繁
	if ttl.Val().Seconds() > float64(s.conf.Captcha.EmailExpire-s.conf.Captcha.EmailInterval) {
		return errors.New("验证码发送过于频繁,请稍后再试")
	}

	// 生成验证码
	captcha := str.RandString(s.conf.Captcha.EmailNumber)

	// 发送邮件
	title := "测试验证码"
	timeStr := strconv.FormatInt(s.conf.Captcha.EmailExpire, 10)
	content := "您本次操作的验证码为:" + captcha + ", 有效期为" + timeStr + "秒"
	err := s.email.SendRegisterMail(email, title, content)
	if err != nil {
		return err
	}

	// 保存验证码到redis缓存中
	err = s.rdb.Set(context.Background(), s.conf.Captcha.CaptchaPrefix+":"+email,
		captcha, time.Duration(s.conf.Captcha.EmailExpire)*time.Second).Err()
	return err
}
