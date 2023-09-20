package config

type Jwt struct {
	Secret                  string `mapstructure:"secret"  yaml:"secret"`
	JwtTtl                  int64  `mapstructure:"jwt_ttl"  yaml:"jwt_ttl"`                                       // token 有效期（秒）
	JwtBlacklistGracePeriod int64  `mapstructure:"jwt_blacklist_grace_period"  yaml:"jwt_blacklist_grace_period"` // 黑名单宽限时间（秒）
	RefreshGracePeriod      int64  `mapstructure:"refresh_grace_period"  yaml:"refresh_grace_period"`             // token 自动刷新宽限时间（秒）
}
