package error

import "errors"

const (
	DEFAULT_ERROR     = 40000 // 默认错误
	VALIDATE_ERROR    = 42200 // 验证错误
	TOKEN_ERROR       = 40100 // Token失效
	FORBIDDEN         = 40300 // 无权限
	NOT_FOUND         = 40400 // 数据不存在
	TOO_MANY_REQUESTS = 42900 // 请求过于频繁
	USER_NOT_FOUND    = 40401 // 用户不存在
	SERVER_ERROR      = 50000 // 服务器错误
	BUSSINESS_ERROR   = 50001 // 业务错误
	PARAM_ERROR       = 50002 // 参数错误
)

var (
	FileNotFoundErr     = errors.New("存储驱动不存在")
	FileNoPermissionErr = errors.New("文件权限不足")
)
