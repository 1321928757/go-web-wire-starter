package config

type Storage struct {
	Default    string  `mapstructure:"default" json:"default" yaml:"default"`                // 默认使用的驱动，local本地 oss阿里云 kodo七牛云
	ImgMaxSize int64   `mapstructure:"img_max_size" json:"img_max_size" yaml:"img_max_size"` // 图片最大上传大小（M）
	Drivers    Drivers `mapstructure:"drivers" json:"drivers" yaml:"drivers"`                //驱动配置
}

type Drivers struct {
	Local     LocalConfig `mapstructure:"local" json:"local" yaml:"local"`
	AliOss    AliConfig   `mapstructure:"ali_oss" json:"ali_oss" yaml:"ali_oss"`
	QiNiu     QiNiuConfig `mapstructure:"qi_niu" json:"qi_niu" yaml:"qi_niu"`
	TecentCos CosConfig   `mapstructure:"tecent_cos" json:"tecent_cos" yaml:"tecent_cos"`
}

// 本地存储配置
type LocalConfig struct {
	RootDir string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	AppUrl  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
}

// 腾讯云cos配置
type CosConfig struct {
	Open      bool   `mapstructure:"open" json:"open" yaml:"open"`
	SecretId  string `mapstructure:"secret_id" json:"secret_id" yaml:"secret_id"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Domain    string `mapstructure:"domain" json:"domain" yaml:"domain"`
	IsSsl     bool   `mapstructure:"is_ssl" json:"is_ssl" yaml:"is_ssl"`
	IsPrivate bool   `mapstructure:"is_private" json:"is_private" yaml:"is_private"`
}

// 阿里云oss配置
type AliConfig struct {
	Open            bool   `mapstructure:"open" json:"open" yaml:"open"`
	AccessKeyId     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret" yaml:"access_key_secret"`
	Bucket          string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	IsSsl           bool   `mapstructure:"is_ssl" json:"is_ssl" yaml:"is_ssl"`
	IsPrivate       bool   `mapstructure:"is_private" json:"is_private" yaml:"is_private"`
}

// 七牛云配置
type QiNiuConfig struct {
	Open      bool   `mapstructure:"open" json:"open" yaml:"open"`
	AccessKey string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Domain    string `mapstructure:"domain" json:"domain" yaml:"domain"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	IsSsl     bool   `mapstructure:"is_ssl" json:"is_ssl" yaml:"is_ssl"`
	IsPrivate bool   `mapstructure:"is_private" json:"is_private" yaml:"is_private"`
}
