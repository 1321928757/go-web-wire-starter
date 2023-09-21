package dao

import (
	"context"
	"encoding/json"
	"go-web-wire-starter/internal/compo/storage"
	"go-web-wire-starter/internal/domain"
	"go-web-wire-starter/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// 文件缓存前缀名
const mediaCacheKeyPre = "media:"

// 规定dao必须实现接口对于的方法
var _ MediaDaoInterface = (*MediaDao)(nil)

type MediaDaoInterface interface {
	Create(ctx context.Context, dm *domain.Media) (*domain.Media, error)
	FindByID(ctx context.Context, id uint64) (*domain.Media, error)
	FindCacheByID(ctx context.Context, id uint64) (*domain.Media, error)
}

type MediaDao struct {
	data    *Data
	log     *zap.Logger
	storage *storage.Storage
}

func NewMediaDao(data *Data, log *zap.Logger, storage *storage.Storage) *MediaDao {
	return &MediaDao{
		data:    data,
		log:     log,
		storage: storage,
	}
}

// 创建媒体文件记录
func (r *MediaDao) Create(ctx context.Context, dm *domain.Media) (*domain.Media, error) {
	var m model.Media

	id, err := r.data.sf.NextID()
	if err != nil {
		return nil, err
	}

	m.ID = id
	m.DiskType = dm.DiskType
	m.SrcType = dm.SrcType
	m.Src = dm.Src

	if err = r.data.DB(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	dm.ID = m.ID

	return dm, nil
}

// 根据ID查找媒体文件记录（不缓存）
func (r *MediaDao) FindByID(ctx context.Context, id uint64) (*domain.Media, error) {
	var m model.Media
	if err := r.data.db.First(&m, id).Error; err != nil {
		return nil, err
	}

	dm := m.ToDomain()
	// 设置访问路径
	dm.SetUrl(r.storage)

	return dm, nil
}

// 根据ID查找媒体文件记录（缓存）
func (r *MediaDao) FindCacheByID(ctx context.Context, id uint64) (*domain.Media, error) {
	cacheKey := mediaCacheKeyPre + strconv.FormatUint(id, 10)

	// 首先从缓存中查找
	exist := r.data.rdb.Exists(ctx, cacheKey).Val()
	if exist == 1 {
		bytes, err := r.data.rdb.Get(ctx, cacheKey).Bytes()
		if err != nil {
			return nil, err
		}
		var media domain.Media
		err = json.Unmarshal(bytes, &media)
		if err != nil {
			return nil, err
		}

		return &media, nil
	}

	// 缓存中没有，从数据库中查找，并存入缓存
	var media model.Media
	err := r.data.db.First(&media, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	dm := media.ToDomain()
	dm.SetUrl(r.storage)
	v, err := json.Marshal(dm)
	if err != nil {
		return nil, err
	}
	r.data.rdb.Set(ctx, cacheKey, v, time.Second*3*24*3600)

	return dm, nil
}
