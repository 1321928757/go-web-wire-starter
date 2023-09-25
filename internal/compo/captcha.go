package compo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/wenlng/go-captcha/captcha"
	"go-web-wire-starter/config"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type CaptchaCompo struct {
	rdb    *redis.Client
	logger *zap.Logger
	config *config.Configuration
}

func NewCaptchaCompo(rdb *redis.Client, logger *zap.Logger, config *config.Configuration) *CaptchaCompo {
	return &CaptchaCompo{rdb: rdb, logger: logger, config: config}
}

// 产生新的点击验证码
func (c *CaptchaCompo) GetNewCaptcha() (dots map[int]captcha.CharDot,
	b64 string, tb64 string, key string, err error) {
	capt := captcha.GetCaptcha()

	// dots为正确点击位置，b64为图形，tb64为略缩图形，key为本次验证码的唯一标识
	dots, b64, tb64, key, err = capt.Generate()

	if err != nil {
		c.logger.Error("生成点击验证码失败:" + err.Error())
		return
	}

	// 序列化value
	data, err := json.Marshal(dots)
	if err != nil {
		c.logger.Error("序列化点击位置时出现错误:" + err.Error())
		return
	}
	// 将验证码位置点击信息存入缓存中
	c.rdb.SetNX(context.Background(), key, data, time.Duration(c.config.Captcha.FigureExpire)*time.Second)
	return
}

// 校验点击验证码
func (c *CaptchaCompo) CheckCaptcha(dotsStr string, key string) (bool, error) {
	// 从缓存中读取
	dotsJson, err := c.rdb.Get(context.Background(), key).Result()
	if err != nil {
		c.logger.Error("从缓存中读取验证码数据出现错误:" + err.Error())
		return false, err
	}

	// 定义map变量来接收反序列化结果
	var dotsRight map[int]captcha.CharDot
	// 反序列化
	err = json.Unmarshal([]byte(dotsJson), &dotsRight)
	if err != nil {
		c.logger.Error("反序列化点击验证码位置信息时出现错误:" + err.Error())
		return false, err
	}

	dots := strings.Split(dotsStr, ",")

	// 校验验证码
	chkRet := false
	if (len(dotsRight) * 2) == len(dots) {
		for i, dot := range dotsRight {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", dots[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", dots[k]), 64)

			// 检测点位置
			// chkRet = captcha.CheckPointDist(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height))

			// 校验点的位置,在原有的区域上添加额外边距进行扩张计算区域,不推荐设置过大的padding
			// 例如：文本的宽和高为30，校验范围x为10-40，y为15-45，此时扩充5像素后校验范围宽和高为40，则校验范围x为5-45，位置y为10-50
			chkRet = captcha.CheckPointDistWithPadding(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy),
				int64(dot.Width), int64(dot.Height), 5)
			if !chkRet {
				break
			}
		}
	}
	return chkRet, nil
}
