package domain

import (
	"go-web-wire-starter/internal/compo/storage"
)

type Media struct {
	ID       uint64 `json:"id"`
	DiskType string `json:"-"`
	SrcType  int8   `json:"-"`
	Src      string `json:"src"`
	Url      string `json:"url"`
}

func (dm *Media) SetUrl(storage *storage.Storage) string {
	url := ""

	// 相对路径
	if dm.SrcType == 1 {
		disk, err := storage.FileDriver(dm.DiskType)
		if err != nil {
			url = ""
		} else {
			url = disk.Url(dm.Src)
		}
	}

	// 外链
	if dm.SrcType == 2 {
		url = dm.Src
	}

	dm.Url = url
	return dm.Url
}
