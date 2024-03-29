package response

import (
	"github.com/gin-gonic/gin"
	cErr "go-web-wire-starter/internal/pkg/error"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		200,
		data,
		"ok",
	})
	c.Abort()
}

// Fail 失败响应
func Fail(c *gin.Context, httpCode int, errorCode int, msg string) {
	c.JSON(httpCode, Response{
		errorCode,
		nil,
		msg,
	})
	c.Abort()
}

// ServerError 服务器错误
func ServerError(c *gin.Context, err interface{}) {
	msg := "Internal Server Error"
	//发布模式下只返回默认错误消息,不泄露内部错误详情
	if gin.Mode() != gin.ReleaseMode {
		if _, ok := err.(error); ok {
			msg = err.(error).Error()
		}
	}
	FailByErr(c, cErr.InternalServer(msg))
}

// FailByErr 错误响应(状态码不为200)
func FailByErr(c *gin.Context, err error) {
	// 判断是否是自定义错误
	v, ok := err.(*cErr.Error)
	if ok {
		// 自定义错误返回自定义错误内容
		Fail(c, v.HttpCode(), v.ErrorCode(), v.Error())
	} else {
		// 非自定义错误返回默认错误
		Fail(c, http.StatusBadRequest, cErr.DEFAULT_ERROR, err.Error())
	}
}

// FailByBussiness 业务错误(状态码为200)
func FailByBussiness(c *gin.Context, msg string) {
	Fail(c, http.StatusOK, cErr.BUSSINESS_ERROR, msg)
}

// FailByParam 参数错误(状态码为200)
func FailByParams(c *gin.Context, msg string) {
	Fail(c, http.StatusOK, cErr.PARAM_ERROR, msg)
}
