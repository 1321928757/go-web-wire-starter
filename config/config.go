package config

type Configuration struct {
	App      App      `mapstructure:"internal" yaml:"internal"`
	Log      Log      `mapstructure:"log" yaml:"log"`
	Database Database `mapstructure:"database" yaml:"database"`
	Jwt      Jwt      `mapstructure:"jwt" yaml:"jwt"`
	Redis    Redis    `mapstructure:"redis" yaml:"redis"`
	Storage  Storage  `mapstructure:"storage" yaml:"storage"`
	Queue    Queue    `mapstructure:"queue" yaml:"queue"`
}
