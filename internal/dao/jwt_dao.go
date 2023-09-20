package dao

import (
	"context"
	"go-web-wire-starter/util/hash"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type JwtDao struct {
	data   *Data
	logger *zap.Logger
}

func NewJwtDao(data *Data, logger *zap.Logger) *JwtDao {
	return &JwtDao{
		data:   data,
		logger: logger,
	}
}

// getBlackListKey 生成标识 JWT 的黑名单键
func (r *JwtDao) getBlackListKey(tokenStr string) string {
	return "jwt_black_list:" + hash.MD5([]byte(tokenStr))
}

// JoinBlackList 将 JWT 添加到黑名单中
func (r *JwtDao) JoinBlackList(ctx context.Context, tokenStr string, joinUnix int64, expires time.Duration) error {
	return r.data.rdb.SetNX(ctx, r.getBlackListKey(tokenStr), joinUnix, expires).Err()
}

// GetBlackJoinUnix 获取 JWT 加入黑名单的时间
func (r *JwtDao) GetBlackJoinUnix(ctx context.Context, tokenStr string) (int64, error) {
	joinUnixStr, err := r.data.rdb.Get(ctx, r.getBlackListKey(tokenStr)).Result()
	if err != nil {
		return 0, err
	}

	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if joinUnixStr == "" || err != nil {
		return 0, err
	}

	return joinUnix, nil
}
