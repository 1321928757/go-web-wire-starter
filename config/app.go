package config

type App struct {
	Env     string `mapstructure:"env" yaml:"env"`
	Port    string `mapstructure:"port" yaml:"port"`
	AppName string `mapstructure:"app_name" yaml:"app_name"`
	AppUrl  string `mapstructure:"app_url" yaml:"app_url"`
}
