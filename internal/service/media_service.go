package service

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/compo/storage"
	"go-web-wire-starter/internal/dao"
	"go-web-wire-starter/internal/domain"
	cErr "go-web-wire-starter/internal/pkg/error"
	"go-web-wire-starter/internal/pkg/request"
	"go.uber.org/zap"
	"path"
)

type MediaService struct {
	conf     *config.Configuration
	log      *zap.Logger
	mediaDao *dao.MediaDao
	storage  *storage.Storage
}

// NewMediaService .
func NewMediaService(conf *config.Configuration, log *zap.Logger, mediaDao *dao.MediaDao,
	s *storage.Storage) *MediaService {
	return &MediaService{conf: conf, log: log, mediaDao: mediaDao, storage: s}
}

// 生成文件夹路径
func (s *MediaService) makeFaceDir(business string) string {
	return s.conf.App.Env + "/" + business
}

// 生成随机文件名
func (s *MediaService) HashName(fileName string) string {
	fileSuffix := path.Ext(fileName)
	return uuid.NewV4().String() + fileSuffix
}

// SaveImage 保存图片（公共读）
func (s *MediaService) SaveImage(ctx *gin.Context, params *request.ImageUpload) (*domain.Media, error) {
	// 读取文件
	file, err := params.Image.Open()
	defer file.Close()
	if err != nil {
		println(11111111111111)
		return nil, cErr.BadRequest("上传失败")
	}

	// 获取存储驱动，生成文件名
	disk, err := s.storage.FileDriver()
	if err != nil {
		return nil, cErr.BadRequest(s.storage.GetDefaultDiskType() + "disk not found")
	}
	localPrefix := ""
	if s.storage.IsLocal() {
		localPrefix = "public" + "/"
	}
	key := s.makeFaceDir(params.Business) + "/" + s.HashName(params.Image.Filename)
	// 上传文件到本地（服务器）
	err = disk.Put(localPrefix+key, file, params.Image.Size)
	if err != nil {
		println(2222222222222)
		return nil, cErr.BadRequest("上传失败")
	}

	// 保存媒体数据到数据库
	m, err := s.mediaDao.Create(ctx, &domain.Media{
		DiskType: s.storage.GetDefaultDiskType(),
		SrcType:  1,
		Src:      key,
		Url:      disk.Url(key),
	})
	if err != nil {
		return nil, cErr.BadRequest("上传失败")
	}

	return m, nil
}

// GetUrlById 根据ID获取文件访问URL
func (s *MediaService) GetUrlById(ctx *gin.Context, id uint64) (string, error) {
	m, err := s.mediaDao.FindCacheByID(ctx, id)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	return m.Url, err
}
