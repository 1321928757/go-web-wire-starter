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

func ServerError(c *gin.Context, err interface{}) {
	msg := "Internal Server Error"
	if gin.Mode() != gin.ReleaseMode {
		if _, ok := err.(error); ok {
			msg = err.(error).Error()
		}
	}
	FailByErr(c, cErr.InternalServer(msg))
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		200,
		data,
		"ok",
	})
	c.Abort()
}

func Fail(c *gin.Context, httpCode int, errorCode int, msg string) {
	c.JSON(httpCode, Response{
		errorCode,
		nil,
		msg,
	})
	c.Abort()
}

func FailByErr(c *gin.Context, err error) {
	v, ok := err.(*cErr.Error)
	if ok {
		Fail(c, v.HttpCode(), v.ErrorCode(), v.Error())
	} else {
		Fail(c, http.StatusBadRequest, cErr.DEFAULT_ERROR, err.Error())
	}
}
