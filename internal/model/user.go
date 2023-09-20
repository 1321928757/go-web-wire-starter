package model

import (
	"go-web-wire-starter/internal/domain"
)

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Name     string `gorm:"size:30;not null;comment:用户名称"`
	Mobile   string `gorm:"size:24;not null;index;comment:用户手机号"`
	Password string `gorm:"not null;default:'';comment:用户密码"`
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
