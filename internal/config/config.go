package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port             string
	PostgresDSN      string
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
	RedisCacheExpiry int
}

func LoadConfig() *Config {
	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		PostgresDSN:      getEnv("POSTGRES_DSN", "user=postgres password=postgres dbname=forecasting sslmode=disable"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          getEnvAsInt("REDIS_DB", 0),
		RedisCacheExpiry: getEnvAsInt("REDIS_CACHE_EXPIRATION", 1),
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}
	return value
}

func getEnvAsInt(name string, defaultVal int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		var value int
		if _, err := fmt.Sscanf(valueStr, "%d", &value); err == nil {
			return value
		}
	}
	return defaultVal
}
