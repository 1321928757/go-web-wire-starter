package domain

import (
	"strconv"
	"time"
)

type User struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Mobile    string    `json:"mobile"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetUid() string {
	return strconv.Itoa(int(u.ID))
}
