package request

import (
	"github.com/go-playground/validator/v10"
	cErr "go-web-wire-starter/internal/pkg/error"
	"regexp"
)

type Validator interface {
	GetMessages() ValidatorMessages
}

type ValidatorMessages map[string]string

var reg = regexp.MustCompile(`\[\d\]`)

// GetError 获取验证错误
func GetError(request interface{}, err error) *cErr.Error {
	if _, isValidatorErrors := err.(validator.ValidationErrors); isValidatorErrors {
		_, isValidator := request.(Validator)

		for _, v := range err.(validator.ValidationErrors) {
			// 若 request 结构体实现 Validator 接口即可实现自定义错误信息
			if isValidator {
				field := v.Field() // 取 request 结构体字段的 'vn' tag 值，未设置 tag 则默认为字段名
				field = reg.ReplaceAllString(field, ".*")
				if message, exist := request.(Validator).GetMessages()[field+"."+v.Tag()]; exist {
					return cErr.ValidateErr(message)
				}
			}
			return cErr.ValidateErr(v.Error())
		}
	}

	return cErr.ValidateErr("Parameter error")
}
