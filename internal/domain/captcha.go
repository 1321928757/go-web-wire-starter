package domain

type ClickCaptcha struct {
	Key         string `json:"key"`          //本次点击验证码的唯一标识
	ImageBase64 string `json:"image_base64"` //本次点击验证码的验证图像
	ThumbBase64 string `json:"thumb_base64"` //本次点击验证码的答案图像
}
