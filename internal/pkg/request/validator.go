package request

import (
	"github.com/go-playground/validator/v10"
)

// Gin 自带验证器返回的错误信息格式不太友好，故在此自定义错误信息
type Validator interface {
	GetMessages() ValidatorMessages
}

type ValidatorMessages map[string]string

// GetErrorMsg 获取错误信息
func GetErrorMsg(request interface{}, err error) string {
	if _, isValidatorErrors := err.(validator.ValidationErrors); isValidatorErrors {
		_, isValidator := request.(Validator)

		for _, v := range err.(validator.ValidationErrors) {
			// 若 request 结构体实现 Validator 接口即可实现自定义错误信息
			if isValidator {
				if message, exist := request.(Validator).GetMessages()[v.Field()+"."+v.Tag()]; exist {
					return message
				}
			}
			return v.Error()
		}
	}

	return "参数错误"
}
