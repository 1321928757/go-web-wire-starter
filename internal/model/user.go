package model

import (
	"go-web-wire-starter/internal/domain"
)

type User struct {
	ID       uint64 `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name;size:30;not null;comment:用户名称"`
	Mobile   string `gorm:"column:mobile;size:24;not null;index;comment:用户手机号"`
	Password string `gorm:"column:password;not null;default:'';comment:用户密码"`
	Timestamps
	SoftDeletes
}

func (m *User) ToDomain() *domain.User {
	return &domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Mobile:    m.Mobile,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (*User) TableName() string {
	return "user"
}
