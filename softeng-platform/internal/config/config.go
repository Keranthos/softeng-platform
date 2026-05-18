package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() *Config {
	// 构建数据库连接字符串 - 使用 softeng_app:123456
	databaseURL := buildDatabaseURL()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: databaseURL,
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func buildDatabaseURL() string {
	// 从环境变量获取配置，如果没有则使用默认值
	user := getEnv("DB_USER", "softeng_app")    // 默认用户
	password := getEnv("DB_PASSWORD", "123456") // 默认密码
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "softeng")

	// 构建 MySQL 连接字符串
	if password == "" {
		return user + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true&loc=Local&charset=utf8mb4"
	}
	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true&loc=Local&charset=utf8mb4"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
