package oss

import (
	"errors"
	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go-web-wire-starter/config"
	"go-web-wire-starter/util/fileUtil"
	"io"
	"strconv"
	"sync"
)

type OssDriver struct {
	config *config.AliConfig
	client *alioss.Client
	bucket *alioss.Bucket
}

var (
	o       *OssDriver
	once    *sync.Once
	initErr error
)

func NewOssDriver(config *config.Configuration) (*OssDriver, error) {
	once = &sync.Once{}
	once.Do(func() {
		o = &OssDriver{}
		o.config = &config.Storage.Drivers.AliOss

		o.client, initErr = alioss.New(o.config.Endpoint, o.config.AccessKeyId, o.config.AccessKeySecret)
		if initErr != nil {
			initErr = errors.New("阿里云oss存储驱动初始化失败：" + initErr.Error())
			return
		}

		o.bucket, initErr = o.client.Bucket(o.config.Bucket)
		if initErr != nil {
			initErr = errors.New("阿里云oss存储驱动初始化失败：" + initErr.Error())
			return
		}
	})
	return o, nil
}

func (o *OssDriver) Put(key string, r io.Reader, dataLength int64) error {
	key = fileUtil.NormalizeKey(key)

	err := o.bucket.PutObject(key, r)
	if err != nil {
		return err
	}

	return nil
}

func (o *OssDriver) PutFile(key string, localFile string) error {
	key = fileUtil.NormalizeKey(key)

	err := o.bucket.PutObjectFromFile(key, localFile)
	if err != nil {
		return err
	}

	return nil
}

func (o *OssDriver) Get(key string) (io.ReadCloser, error) {
	key = fileUtil.NormalizeKey(key)

	body, err := o.bucket.GetObject(key)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (o *OssDriver) Rename(srcKey string, destKey string) error {
	srcKey = fileUtil.NormalizeKey(srcKey)
	destKey = fileUtil.NormalizeKey(destKey)

	_, err := o.bucket.CopyObject(srcKey, destKey)
	if err != nil {
		return err
	}

	err = o.Delete(srcKey)
	if err != nil {
		return err
	}

	return nil
}

func (o *OssDriver) Copy(srcKey string, destKey string) error {
	srcKey = fileUtil.NormalizeKey(srcKey)
	destKey = fileUtil.NormalizeKey(destKey)

	_, err := o.bucket.CopyObject(srcKey, destKey)
	if err != nil {
		return err
	}

	return nil
}

func (o *OssDriver) Exists(key string) (bool, error) {
	key = fileUtil.NormalizeKey(key)

	return o.bucket.IsObjectExist(key)
}

func (o *OssDriver) Size(key string) (int64, error) {
	key = fileUtil.NormalizeKey(key)

	props, err := o.bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return 0, err
	}

	size, err := strconv.ParseInt(props.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func (o *OssDriver) Delete(key string) error {
	key = fileUtil.NormalizeKey(key)

	err := o.bucket.DeleteObject(key)
	if err != nil {
		return err
	}

	return nil
}

func (o *OssDriver) Url(key string) string {
	var prefix string
	key = fileUtil.NormalizeKey(key)

	if o.config.IsSsl {
		prefix = "https://"
	} else {
		prefix = "http://"
	}

	// 如果是私有存储桶，则通过 o.bucket.SignURL 方法生成带有签名的 URL
	if o.config.IsPrivate {
		url, err := o.bucket.SignURL(key, alioss.HTTPGet, 3600)
		if err == nil {
			return url
		}
	}

	return prefix + o.config.Bucket + "." + o.config.Endpoint + "/" + key
}
