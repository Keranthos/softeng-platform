package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	// CORS：逗号分隔的 Origin 白名单（环境变量 CORS_ALLOWED_ORIGINS）；空则用本地开发默认
	CORSAllowedOrigins []string
	// RAG：Gemini（与 blog/agent 一致，优先于 OpenAI）
	GoogleAPIKey      string  // GOOGLE_API_KEY 或 GEMINI_API_KEY
	GeminiAPIBase     string  // GEMINI_API_BASE，默认 Google 官方
	GeminiProxyToken  string  // 走自定义反代时 GEMINI_PROXY_TOKEN
	GeminiChatModel   string  // GEMINI_CHAT_MODEL
	GeminiTemperature float64 // GEMINI_TEMPERATURE
	// RAG：OpenAI 兼容（次选；两者皆空则 demo）
	OpenAIAPIKey  string
	OpenAIBaseURL string
	OpenAIModel   string
}

func LoadConfig() *Config {
	// 构建数据库连接字符串 - 使用 softeng_app:123456
	databaseURL := buildDatabaseURL()

	gKey := strings.TrimSpace(getEnv("GOOGLE_API_KEY", ""))
	if gKey == "" {
		gKey = strings.TrimSpace(getEnv("GEMINI_API_KEY", ""))
	}
	gTemp := 0.2
	if v := strings.TrimSpace(os.Getenv("GEMINI_TEMPERATURE")); v != "" {
		if x, err := strconv.ParseFloat(v, 64); err == nil {
			gTemp = x
		}
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        databaseURL,
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		CORSAllowedOrigins: parseCORSAllowedOrigins(getEnv("CORS_ALLOWED_ORIGINS", "")),
		GoogleAPIKey:      gKey,
		GeminiAPIBase:     getEnv("GEMINI_API_BASE", "https://generativelanguage.googleapis.com"),
		GeminiProxyToken:  getEnv("GEMINI_PROXY_TOKEN", ""),
		GeminiChatModel:   getEnv("GEMINI_CHAT_MODEL", "gemini-2.5-flash"),
		GeminiTemperature: gTemp,
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
		OpenAIBaseURL:     getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OpenAIModel:       getEnv("OPENAI_MODEL", "gpt-4o-mini"),
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

func defaultCORSOrigins() []string {
	return []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3001",
	}
}

// parseCORSAllowedOrigins 解析 CORS_ALLOWED_ORIGINS；空或全空白则返回本地开发默认
func parseCORSAllowedOrigins(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultCORSOrigins()
	}
	parts := strings.Split(raw, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return defaultCORSOrigins()
	}
	return out
}
