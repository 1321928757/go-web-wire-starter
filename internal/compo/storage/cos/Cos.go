/*
腾讯云Cos存储驱动
*/
package cos

import (
	"context"
	tecentcos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"go-web-wire-starter/config"
	"go-web-wire-starter/util/fileUtil"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// local 本地存储驱动
type CosDriver struct {
	client *tecentcos.Client
	config *config.CosConfig
}

var (
	c    *CosDriver
	once *sync.Once
)

// Init 初始化本地存储驱动
func NewCosDriver(config *config.Configuration) (*CosDriver, error) {
	once = &sync.Once{}
	once.Do(func() { //sync.Once 确保其中的函数只会执行一次，防止多次调用init导致多次注册引擎
		cosConfig := &config.Storage.Drivers.TecentCos
		//初始化cos驱动
		baseUrl := "https://" + cosConfig.Bucket + ".cos." + cosConfig.Domain + ".myqcloud.com"
		u, _ := url.Parse(baseUrl)
		b := &tecentcos.BaseURL{BucketURL: u}
		c = &CosDriver{
			client: tecentcos.NewClient(b, &http.Client{
				Transport: &tecentcos.AuthorizationTransport{
					SecretID:  cosConfig.SecretId,
					SecretKey: cosConfig.SecretKey,
					Transport: &debug.DebugRequestTransport{
						RequestHeader:  true,
						RequestBody:    false,
						ResponseHeader: true,
						ResponseBody:   false,
					},
				},
			}),
			config: cosConfig,
		}
	})
	return c, nil
}

func (c *CosDriver) Put(key string, r io.Reader, dataLength int64) error {
	//腾讯云cos上传文件
	//打印key
	_, err := c.client.Object.Put(context.Background(), key, r, nil)
	if err != nil {
		err.Error()
	}

	return err
}

func (c *CosDriver) PutFile(key string, localFile string) error {
	_, err := c.client.Object.PutFromFile(context.Background(), key, localFile, nil)
	if err != nil {
		err.Error()
	}
	return err
}

func (c *CosDriver) Get(key string) (io.ReadCloser, error) {
	response, err := c.client.Object.Get(context.Background(), key, nil)
	file := response.Body
	if err != nil {
		err.Error()
	}
	return file, err
}

func (c *CosDriver) Rename(srcKey string, destKey string) error {
	srcKey = fileUtil.NormalizeKey(srcKey)
	destKey = fileUtil.NormalizeKey(destKey)

	// 改名其实就是复制文件，然后删除原文件
	_, _, err := c.client.Object.Copy(context.Background(), srcKey, destKey, nil)
	if err != nil {
		return err
	}

	_, err = c.client.Object.Delete(context.Background(), srcKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *CosDriver) Copy(srcKey string, destKey string) error {
	srcKey = fileUtil.NormalizeKey(srcKey)
	destKey = fileUtil.NormalizeKey(destKey)

	_, _, err := c.client.Object.Copy(context.Background(), srcKey, destKey, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *CosDriver) Exists(key string) (bool, error) {
	key = fileUtil.NormalizeKey(key)
	return c.client.Object.IsExist(context.Background(), key)
}

func (c *CosDriver) Size(key string) (int64, error) {
	key = fileUtil.NormalizeKey(key)

	// 获取文件信息
	head, err := c.client.Object.Head(context.Background(), key, nil)
	if err != nil {
		return 0, err
	}

	return head.ContentLength, nil
}

func (c *CosDriver) Delete(key string) error {
	key = fileUtil.NormalizeKey(key)

	_, err := c.client.Object.Delete(context.Background(), key)
	if err != nil {
		return err
	}

	return nil
}

func (c *CosDriver) Url(key string) string {
	var prefix string
	key = fileUtil.NormalizeKey(key)

	if c.config.IsSsl {
		prefix = "https://"
	} else {
		prefix = "http://"
	}

	// 如果是私有存储桶，则通过 o.bucket.SignURL 方法生成带有签名的 URL
	if c.config.IsPrivate {
		presignedURL, err := c.client.Object.GetPresignedURL(context.Background(), http.MethodGet, key,
			c.config.SecretId, c.config.SecretKey, time.Hour, nil)
		if err == nil {
			return presignedURL.String()
		}
	}

	return prefix + c.config.Bucket + ".cos." + c.config.Domain + ".myqcloud.com" + "/" + key
}
