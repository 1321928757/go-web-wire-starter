package fileUtil

import (
	"errors"
	"os"
	"strings"
)

var (
	FileNotFoundErr     = errors.New("fileDriver not found")
	FileNoPermissionErr = errors.New("permission denied")
)

// NormalizeKey 字符串进行标准化处理.字符串中的反斜杠 \ 替换为正斜杠 /，去除字符串中的空格，过滤字符串中的换行符和一些特殊字符。
func NormalizeKey(key string) string {
	key = strings.Replace(key, "\\", "/", -1)
	key = strings.Replace(key, " ", "", -1)
	key = filterNewLines(key)

	return key
}

// filterNewLines 过滤字符串中的换行符和一些特殊字符。
func filterNewLines(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case 0x000A, 0x000B, 0x000C, 0x000D, 0x0085, 0x2028, 0x2029:
			return -1
		default:
			return r
		}
	}, s)
}

// OpenAsReadOnly 以只读方式打开指定的文件
func OpenAsReadOnly(key string) (*os.File, os.FileInfo, error) {
	fd, err := os.Open(key)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, FileNotFoundErr
		}
		if os.IsPermission(err) {
			return nil, nil, FileNoPermissionErr
		}
		return nil, nil, err
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, nil, err
	}

	return fd, stat, nil
}

// IsValidImageType 校验文件是否为图片类型， contentType 为文件类型（格式符合请求头的Content-Type）
func IsValidImageType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif":
		return true
	default:
		return false
	}
}
