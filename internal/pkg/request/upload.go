package request

import "mime/multipart"

type ImageUpload struct {
	Business string                `form:"business" binding:"required"`
	Image    *multipart.FileHeader `form:"image" binding:"required"`
}

func (imageUpload ImageUpload) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"Business.required": "业务类型不能为空",
		"Image.required":    "请选择图片",
	}
}
