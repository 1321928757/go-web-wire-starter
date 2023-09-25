package request

type CaptchaClickParams struct {
	Key  string `form:"key" binding:"required"`
	Dots string `form:"dots" binding:"required"`
}

func (captchaClickParams CaptchaClickParams) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"Key.required":  "验证码key不能为空",
		"Dots.required": "用户点击位置信息不能为空",
	}
}
