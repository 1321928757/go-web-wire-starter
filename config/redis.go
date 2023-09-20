package config

type Redis struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	DB       int    `mapstructure:"db" yaml:"db"`
	Password string `mapstructure:"password" yaml:"password"`
}
