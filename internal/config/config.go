package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 用于集中管理应用配置。
// 学习阶段通过环境变量读取，避免引入额外配置库，便于理解。
type Config struct {
	App struct {
		Name      string
		Port      string
		APIPrefix string
	}
	Database struct {
		Host      string
		Port      string
		User      string
		Password  string
		Name      string
		Charset   string
		ParseTime string
		Loc       string
	}
	JWT struct {
		Secret      string
		ExpireHours int
	}
}

// Load 从环境变量中加载配置并设置默认值。
func Load() *Config {
	cfg := &Config{}

	cfg.App.Name = getEnv("APP_NAME", "mini-mall-gin-gorm")
	cfg.App.Port = getEnv("APP_PORT", "8080")
	cfg.App.APIPrefix = getEnv("APP_API_PREFIX", "/api")

	cfg.Database.Host = getEnv("DB_HOST", "127.0.0.1")
	cfg.Database.Port = getEnv("DB_PORT", "3306")
	cfg.Database.User = getEnv("DB_USER", "root")
	cfg.Database.Password = getEnv("DB_PASSWORD", "123456")
	cfg.Database.Name = getEnv("DB_NAME", "mini_mall")
	cfg.Database.Charset = getEnv("DB_CHARSET", "utf8mb4")
	cfg.Database.ParseTime = getEnv("DB_PARSE_TIME", "True")
	cfg.Database.Loc = getEnv("DB_LOC", "Local")

	cfg.JWT.Secret = getEnv("JWT_SECRET", "replace-with-your-secret")
	cfg.JWT.ExpireHours = getEnvInt("JWT_EXPIRE_HOURS", 24)

	return cfg
}

// DSN 生成 MySQL 连接字符串。
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.Charset,
		c.Database.ParseTime,
		c.Database.Loc,
	)
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
