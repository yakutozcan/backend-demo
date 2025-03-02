package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port             string `mapstructure:"port" yaml:"port"`
	DatabaseUser     string `mapstructure:"database_user" yaml:"database_user"`
	DatabasePassword string `mapstructure:"database_password" yaml:"database_password"`
	DatabaseHost     string `mapstructure:"database_host" yaml:"database_host"`
	DatabasePort     string `mapstructure:"database_port" yaml:"database_port"`
	DatabaseName     string `mapstructure:"database_name" yaml:"database_name"`
	DatabaseSSLMode  string `mapstructure:"database_sslmode" yaml:"database_sslmode"`
	DatabaseTimezone string `mapstructure:"database_timezone" yaml:"database_timezone"`
	RedisHost        string `mapstructure:"redis_host" yaml:"redis_host"`
	RedisPort        string `mapstructure:"redis_port" yaml:"redis_port"`
	RedisPassword    string `mapstructure:"redis_password" yaml:"redis_password"`
	RedisDB          int    `mapstructure:"redis_db" yaml:"redis_db"`
}

func Read() *AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PWD/config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	return &appConfig
}
