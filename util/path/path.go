package path

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// RootPath 获取项目根目录绝对路径
func RootPath() string {
	var rootDir string

	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 排除磁盘卷标的情况
	//根路径可能会被表示为 "D:" 或类似的磁盘卷标，而不是具体的文件系统路径。
	//这种情况下，可以考虑将 "D:" 排除在外，只返回实际的文件系统路径。
	if strings.Contains(exePath, ":\\") {
		return exePath
	}

	rootDir = filepath.Dir(filepath.Dir(exePath))

	tmpDir := os.TempDir()
	if strings.Contains(exePath, tmpDir) {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			rootDir = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
		}
	}

	return rootDir
}

// Exists 路径是否存在
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
