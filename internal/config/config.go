package config

import (
	"os"
)

// AppConfig хранит конфигурацию приложения
type AppConfig struct {
	ServerPort string
	JWTSecret  []byte
	DataDir    string
	StaticDir  string
	GinMode    string
}

// App — глобальная конфигурация
var App AppConfig

// Init инициализирует конфигурацию из переменных окружения
func Init() {
	App = AppConfig{
		ServerPort: getEnv("PORT", "8080"),
		JWTSecret:  []byte(getEnv("JWT_SECRET", "your-secret-key-change-in-production-min-32-chars!")),
		DataDir:    getEnv("DATA_DIR", "data"),
		StaticDir:  getEnv("STATIC_DIR", "static"),
		GinMode:    getEnv("GIN_MODE", "debug"),
	}
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
