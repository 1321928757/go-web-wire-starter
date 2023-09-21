package kodo

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	qiniuStorage "github.com/qiniu/go-sdk/v7/storage"
	"go-web-wire-starter/config"
	"go-web-wire-starter/util/fileUtil"
	"io"
	"net/http"
	"sync"
	"time"
)

type KodoDriver struct {
	config        *config.QiNiuConfig
	putPolicy     *qiniuStorage.PutPolicy
	mac           *qbox.Mac
	formUploader  *qiniuStorage.FormUploader
	bucketManager *qiniuStorage.BucketManager
}

var (
	k    *KodoDriver
	once *sync.Once
)

func NewKodoDriver(config *config.Configuration) (*KodoDriver, error) {
	once = &sync.Once{}
	once.Do(func() {
		k = &KodoDriver{}
		k.config = &config.Storage.Drivers.QiNiu

		k.putPolicy = &qiniuStorage.PutPolicy{
			Scope: k.config.Bucket,
		}
		k.mac = qbox.NewMac(k.config.AccessKey, k.config.SecretKey)

		cfg := qiniuStorage.Config{
			UseHTTPS:      k.config.IsSsl,
			UseCdnDomains: false,
		}

		k.formUploader = qiniuStorage.NewFormUploader(&cfg)
		k.bucketManager = qiniuStorage.NewBucketManager(k.mac, &cfg)
	})
	return k, nil
}

func (k *KodoDriver) Put(key string, r io.Reader, dataLength int64) error {
	key = fileUtil.NormalizeKey(key)

	upToken := k.putPolicy.UploadToken(k.mac)
	ret := qiniuStorage.PutRet{}
	err := k.formUploader.Put(context.Background(), &ret, upToken, key, r, dataLength, nil)
	if err != nil {
		return err
	}

	return nil
}

func (k *KodoDriver) PutFile(key string, localFile string) error {
	key = fileUtil.NormalizeKey(key)

	upToken := k.putPolicy.UploadToken(k.mac)
	ret := qiniuStorage.PutRet{}
	err := k.formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, nil)
	if err != nil {
		return err
	}

	return nil
}

func (k *KodoDriver) Get(key string) (io.ReadCloser, error) {
	key = fileUtil.NormalizeKey(key)

	resp, err := http.Get(k.Url(key))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (k *KodoDriver) Rename(srcKey string, destKey string) error {
	srcKey = fileUtil.NormalizeKey(srcKey)
	destKey = fileUtil.NormalizeKey(destKey)

	err := k.bucketManager.Move(k.config.Bucket, srcKey, k.config.Bucket, destKey, true)
	if err != nil {
		return err
	}

	return nil
}

func (k *KodoDriver) Copy(srcKey string, destKey string) error {
	srcKey = fileUtil.NormalizeKey(srcKey)
	destKey = fileUtil.NormalizeKey(destKey)

	err := k.bucketManager.Copy(k.config.Bucket, srcKey, k.config.Bucket, destKey, true)
	if err != nil {
		return err
	}

	return nil
}

func (k *KodoDriver) Exists(key string) (bool, error) {
	key = fileUtil.NormalizeKey(key)

	_, err := k.bucketManager.Stat(k.config.Bucket, key)
	if err != nil {
		if err.Error() == "no such file or directory" {
			err = nil
		}
		return false, err
	}

	return true, nil
}

func (k *KodoDriver) Size(key string) (int64, error) {
	key = fileUtil.NormalizeKey(key)

	fileInfo, err := k.bucketManager.Stat(k.config.Bucket, key)
	if err != nil {
		return 0, err
	}

	return fileInfo.Fsize, nil
}

func (k *KodoDriver) Delete(key string) error {
	key = fileUtil.NormalizeKey(key)

	err := k.bucketManager.Delete(k.config.Bucket, key)
	if err != nil {
		return err
	}

	return nil
}

func (k *KodoDriver) Url(key string) string {
	var prefix string

	key = fileUtil.NormalizeKey(key)

	if k.config.IsSsl {
		prefix = "https://"
	} else {
		prefix = "http://"
	}

	if k.config.IsPrivate {
		deadline := time.Now().Add(time.Second * 3600).Unix() // 1小时有效期
		return prefix + qiniuStorage.MakePrivateURL(k.mac, k.config.Domain, key, deadline)
	}

	return prefix + qiniuStorage.MakePublicURL(k.config.Domain, key)
}
