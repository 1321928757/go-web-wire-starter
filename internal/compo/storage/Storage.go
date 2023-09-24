package storage

import (
	"fmt"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/compo/storage/cos"
	"go-web-wire-starter/internal/compo/storage/kodo"
	"go-web-wire-starter/internal/compo/storage/local"
	"go-web-wire-starter/internal/compo/storage/oss"
	"go.uber.org/zap"
	"io"
)

type DriverName string

// 存储驱动名
const (
	Local DriverName = "local"      // 本地
	KoDo  DriverName = "qi_niu"     // 七牛云
	Oss   DriverName = "ali_oss"    // 阿里云
	Cos   DriverName = "tecent_cos" // 腾讯云
)

// StorageDriver 存储驱动接口，自定义的存储驱动需要实现该接口
type StorageDriver interface {
	// Put 上传文件
	Put(key string, r io.Reader, dataLength int64) error
	// PutFile 上传本地文件
	PutFile(key string, localFile string) error
	// Get 获取文件
	Get(key string) (io.ReadCloser, error)
	// Rename 重命名文件
	Rename(srcKey string, destKey string) error
	// Copy 复制文件
	Copy(srcKey string, destKey string) error
	// Exists 判断文件是否存在
	Exists(key string) (bool, error)
	// Size 获取文件大小
	Size(key string) (int64, error)
	// Delete 删除文件
	Delete(key string) error
	// Url 获取文件访问URL
	Url(key string) string
}

type Storage struct {
	logger  *zap.Logger
	conf    *config.Configuration
	drivers map[DriverName]StorageDriver
}

// NewStorage 初始化存储驱动
func NewStorage(conf *config.Configuration, logger *zap.Logger, cos *cos.CosDriver,
	local *local.LocalDriver, oss *oss.OssDriver, kodo *kodo.KodoDriver) *Storage {
	storage := &Storage{
		logger:  logger,
		conf:    conf,
		drivers: make(map[DriverName]StorageDriver),
	}
	// 注册存储驱动
	storage.Register(Local, local) // 本地
	storage.Register(Cos, cos)     // 腾讯云cos
	storage.Register(Oss, oss)     // 阿里云oss
	storage.Register(KoDo, kodo)   // 七牛云kodo

	return storage
}

// 判断是否为本地存储
func (s *Storage) IsLocal() bool {
	return s.conf.Storage.Default == string(Local)
}

// 获取默认存储驱动名
func (s *Storage) GetDefaultDiskType() string {
	return s.conf.Storage.Default
}

// 注册存储驱动
func (storage *Storage) Register(name DriverName, disk StorageDriver) {
	if disk == nil {
		panic("storage: Register disk is nil")
	}
	storage.drivers[name] = disk
}

// 根据驱动名获取对于存储驱动
func (storage *Storage) GetDriverByName(name DriverName) (StorageDriver, error) {
	disk, exist := storage.drivers[name]
	//遍历打印所有驱动
	if !exist {
		return nil, fmt.Errorf("获取存储驱动错误: 未找到对于存储驱动 %q", name)
	}
	return disk, nil
}

// 使用过程中获取存储驱动实例的统一入口
func (storage *Storage) FileDriver(disk ...string) (StorageDriver, error) {
	// 若未传参，默认使用配置文件驱动
	// 读取默认配置
	diskNameStr := storage.conf.Storage.Default
	// 将str字符转换为自定义的DriverName类型
	diskName := DriverName(diskNameStr)
	// 若传参，使用传参的驱动名
	if len(disk) > 0 {
		diskName = DriverName(disk[0])
	}
	s, err := storage.GetDriverByName(diskName)
	if err != nil {
		return nil, err
	}
	return s, nil
}
