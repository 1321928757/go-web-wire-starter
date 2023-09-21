package model

import "go-web-wire-starter/internal/domain"

type Media struct {
	ID       uint64 `gorm:"column:id;primaryKey"`
	DiskType string `gorm:"column:disk_type;size:20;index;not null;comment:存储类型"`
	SrcType  int8   `gorm:"column:src_type;not null;comment:链接类型 1相对路径 2外链"`
	Src      string `gorm:"column:src;not null;comment:资源链接"`
	Timestamps
}

func (m *Media) ToDomain() *domain.Media {
	return &domain.Media{
		ID:       m.ID,
		DiskType: m.DiskType,
		SrcType:  m.SrcType,
		Src:      m.Src,
	}
}

func (*Media) TableName() string {
	return "media"
}
