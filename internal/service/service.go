package service

import (
	"context"
	"github.com/google/wire"
)

// ProviderSet Provider对象集合
var ProviderSet = wire.NewSet(NewUserService, NewJwtService)

// Transaction 新增事务接口方法
type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}
