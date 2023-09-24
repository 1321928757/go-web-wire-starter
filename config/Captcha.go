package config

type Captcha struct {
	// 验证码缓存key前缀
	CaptchaPrefix string `mapstructure:"captcha_prefix" yaml:"captcha_prefix"`
	// 邮箱验证码长度
	EmailNumber int `mapstructure:"email_number" yaml:"email_number"`
	// 图形验证码长度
	FigureNumber int `mapstructure:"figure_number" yaml:"figure_number"`
	// 邮箱验证码生命周期
	EmailExpire int64 `mapstructure:"email_expire" yaml:"email_timeout"`
	// 图形验证码生命周期
	FigureExpire int64 `mapstructure:"figure_expire" yaml:"figure_timeout"`
	// 发送邮箱验证码最小间隔时长
	EmailInterval int64 `mapstructure:"email_interval" yaml:"email_interval"`
}
