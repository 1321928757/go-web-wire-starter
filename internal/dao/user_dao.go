package dao

import (
	"context"
	"go-web-wire-starter/internal/domain"
	"go-web-wire-starter/internal/model"
	"go.uber.org/zap"
)

// 规定dao必须实现接口对于的方法
var _ UserDaoInterface = (*UserDao)(nil)

type UserDaoInterface interface {
	FindByID(context.Context, uint64) (*domain.User, error)
	FindByMobile(context.Context, string) (*domain.User, error)
	Create(context.Context, *domain.User) (*domain.User, error)
}

type UserDao struct {
	data   *Data
	logger *zap.Logger
}

func NewUserDao(data *Data, logger *zap.Logger) *UserDao {
	return &UserDao{
		data:   data,
		logger: logger,
	}
}

func (r *UserDao) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	var user model.User
	if err := r.data.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *UserDao) FindByMobile(ctx context.Context, mobile string) (*domain.User, error) {
	var user model.User

	if err := r.data.db.Where(&domain.User{Mobile: mobile}).First(&user).Error; err != nil {
		return nil, err
	}

	return user.ToDomain(), nil
}

func (r *UserDao) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	var user model.User

	id, err := r.data.sf.NextID()
	if err != nil {
		return nil, err
	}
	user.ID = id
	user.Name = u.Name
	user.Mobile = u.Mobile
	user.Password = u.Password

	if err = r.data.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user.ToDomain(), nil
}
