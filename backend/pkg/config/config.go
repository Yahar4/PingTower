package config

import "github.com/spf13/viper"

type (
	Config struct {
		Database *DatabaseConfig `mapstructure:"database"`
	}

	DatabaseConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"database"`
		SSLMode  string `mapstructure:"sslmode"`
	}
)

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
