package config

type Database struct {
	Driver              string `mapstructure:"driver" yaml:"driver"`
	Host                string `mapstructure:"host" yaml:"host"`
	Port                int    `mapstructure:"port" yaml:"port"`
	Database            string `mapstructure:"database"  yaml:"database"`
	TablePrefix         string `mapstructure:"table_prefix"  yaml:"table_prefix"`
	UserName            string `mapstructure:"username"  yaml:"username"`
	Password            string `mapstructure:"password"  yaml:"password"`
	Charset             string `mapstructure:"charset"  yaml:"charset"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns"  yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns"  yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode"  yaml:"log_mode"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer"  yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename"  yaml:"log_filename"`
}
