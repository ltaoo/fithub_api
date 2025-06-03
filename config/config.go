package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 存储所有配置
type Config struct {
	// 服务器配置
	ServerAddress string
	Environment   string
	LogLevel      string

	// 数据库配置
	DBType     string // mysql, postgres, sqlite
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPath     string // 用于SQLite

	// 迁移配置
	MigrationsPath string
}

// LoadConfig 从环境变量或配置文件加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// 设置默认值
	viper.SetDefault("SERVER_ADDRESS", ":8080")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DB_TYPE", "sqlite")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "myapi")
	viper.SetDefault("DB_PATH", "./myapi.db")
	viper.SetDefault("MIGRATIONS_PATH", "file:///migrations")

	config := &Config{
		ServerAddress:  viper.GetString("SERVER_ADDRESS"),
		Environment:    viper.GetString("ENVIRONMENT"),
		LogLevel:       viper.GetString("LOG_LEVEL"),
		DBType:         viper.GetString("DB_TYPE"),
		DBHost:         viper.GetString("DB_HOST"),
		DBPort:         viper.GetString("DB_PORT"),
		DBUser:         viper.GetString("DB_USER"),
		DBPassword:     viper.GetString("DB_PASSWORD"),
		DBName:         viper.GetString("DB_NAME"),
		DBPath:         viper.GetString("DB_PATH"),
		MigrationsPath: viper.GetString("MIGRATIONS_PATH"),
	}

	return config, nil
}
