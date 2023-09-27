package config

type Limiter struct {
	Capacity int   `mapstructure:"capacity"  yaml:"capacity"`
	Rate     int64 `mapstructure:"rate"  yaml:"rate"` // token 有效期（秒）
}
