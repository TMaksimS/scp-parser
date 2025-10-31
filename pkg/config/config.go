package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type SCPConfig struct {
	URL       string
	UserAgent string
	Timeout   int
}

type PGConfig struct {
	DBUser     string
	DBPass     string
	DBHost     string
	DBPort     string
	DBName     string
	MaxAttemps int
}

type Config struct {
	SCP SCPConfig
	DB  PGConfig
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("Warning: .env file not found: %v", err)
	}
	return &Config{
		SCP: SCPConfig{
			URL:       getEnv("URL", ""),
			UserAgent: getEnv("UserAgent", ""),
			Timeout:   getEnvInt("ReqTimeout", 5),
		},
		DB: PGConfig{
			DBUser:     getEnv("DBUser", ""),
			DBPass:     getEnv("DBPass", ""),
			DBHost:     getEnv("DBHost", ""),
			DBPort:     getEnv("DBPort", ""),
			DBName:     getEnv("DBName", ""),
			MaxAttemps: getEnvInt("MaxAttemps", 5),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		res, err := strconv.Atoi(value)
		if err != nil {
			fmt.Errorf("Error when converting %s to int", key)
		}
		return res
	}
	return defaultValue
}
