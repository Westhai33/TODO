package config

import (
	"log"
	"os"
)

// Config представляет структуру для конфигурации базы данных
type Config struct {
	DBUser     string // Пользователь базы данных
	DBPassword string // Пароль базы данных
	DBName     string // Имя базы данных
	DBHost     string // Хост базы данных
	DBPort     string // Порт базы данных
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "postgres")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")

	log.Printf("Настройки базы данных: user=%s, dbname=%s, host=%s, port=%s", dbUser, dbName, dbHost, dbPort)

	return &Config{
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBHost:     dbHost,
		DBPort:     dbPort,
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
