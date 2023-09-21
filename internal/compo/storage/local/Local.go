package local

import (
	"go-web-wire-starter/config"
	commonError "go-web-wire-starter/internal/pkg/error"
	"go-web-wire-starter/util/fileUtil"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// local 本地存储驱动
type LocalDriver struct {
	config *config.LocalConfig
}

var (
	l    *LocalDriver
	once *sync.Once
)

func NewLocalDriver(config *config.Configuration) (*LocalDriver, error) {
	once = &sync.Once{}
	once.Do(func() {
		l = &LocalDriver{
			config: &config.Storage.Drivers.Local,
		}
	})
	return l, nil
}

func (l *LocalDriver) getPath(key string) string {
	key = fileUtil.NormalizeKey(key)
	return filepath.Join(l.config.RootDir, key)
}

func (l *LocalDriver) Put(key string, r io.Reader, dataLength int64) error {
	// 处理文件路径
	path := l.getPath(key)
	// 提取文件夹路径，如果文件夹不存在则创建
	dir, _ := filepath.Split(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// 创建文件
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		if os.IsPermission(err) {
			return commonError.FileNoPermissionErr
		}
		return err
	}
	defer fd.Close()

	// 写入文件
	_, err = io.Copy(fd, r)

	return err
}

func (l *LocalDriver) PutFile(key string, localFile string) error {
	path := l.getPath(localFile)

	fd, fileInfo, err := fileUtil.OpenAsReadOnly(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	return l.Put(key, fd, fileInfo.Size())
}

func (l *LocalDriver) Get(key string) (io.ReadCloser, error) {
	path := l.getPath(key)

	fd, _, err := fileUtil.OpenAsReadOnly(path)
	if err != nil {
		return nil, err
	}

	return fd, nil
}

func (l *LocalDriver) Rename(srcKey string, destKey string) error {
	srcPath := l.getPath(srcKey)
	ok, err := l.Exists(srcPath)
	if err != nil {
		return err
	}
	if !ok {
		return commonError.FileNoPermissionErr
	}

	destPath := l.getPath(destKey)
	dir, _ := filepath.Split(destPath)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return os.Rename(srcPath, destPath)
}

func (l *LocalDriver) Copy(srcKey string, destKey string) error {
	srcPath := l.getPath(srcKey)
	srcFd, _, err := fileUtil.OpenAsReadOnly(srcPath)
	if err != nil {
		return err
	}
	defer srcFd.Close()

	destPath := l.getPath(destKey)
	dir, _ := filepath.Split(destPath)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	destFd, err := os.OpenFile(destPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		if os.IsPermission(err) {
			return commonError.FileNoPermissionErr
		}
		return err
	}
	defer destFd.Close()

	_, err = io.Copy(destFd, srcFd)
	if err != nil {
		return err
	}

	return nil
}

func (l *LocalDriver) Exists(key string) (bool, error) {
	path := l.getPath(key)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		if os.IsPermission(err) {
			return false, commonError.FileNoPermissionErr
		}
		return false, err
	}

	return true, nil
}

func (l *LocalDriver) Size(key string) (int64, error) {
	path := l.getPath(key)
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, commonError.FileNotFoundErr
		}
		if os.IsPermission(err) {
			return 0, commonError.FileNoPermissionErr
		}
		return 0, err
	}

	return fileInfo.Size(), nil
}

func (l *LocalDriver) Delete(key string) error {
	path := l.getPath(key)
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return commonError.FileNotFoundErr
		}
		if os.IsPermission(err) {
			return commonError.FileNoPermissionErr
		}
		return err
	}

	return nil
}

// Url 获取文件访问链接
func (l *LocalDriver) Url(key string) string {
	return l.config.AppUrl + "/" + fileUtil.NormalizeKey(key)
}
