package config

import "time"

type Email struct {
	SenderName     string        `mapstructure:"sender_name" yaml:"sender_name"`
	SenderEmail    string        `mapstructure:"sender_email" yaml:"sender_email"`
	SenderPassword string        `mapstructure:"sender_password" yaml:"sender_password"`
	Host           string        `mapstructure:"host" yaml:"host"`
	Port           string        `mapstructure:"port" yaml:"port"`
	MaxConnection  int           `mapstructure:"max_connection" yaml:"max_connection"`
	MaxTimeout     time.Duration `mapstructure:"max_timeout" yaml:"max_timeout"`
}
