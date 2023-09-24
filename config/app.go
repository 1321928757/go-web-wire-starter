package config

type App struct {
	// 当前开发环境
	Env string `mapstructure:"env" yaml:"env"`
	// 服务端口
	Port string `mapstructure:"port" yaml:"port"`
	// 服务名称
	AppName string `mapstructure:"app_name" yaml:"app_name"`
	// 服务地址
	AppUrl string `mapstructure:"app_url" yaml:"app_url"`
}
