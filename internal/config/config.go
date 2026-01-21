package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	AllowedOrigin string
	DBPath        string
	AdminToken    string
	AdminPassword string
}

func Load() Config {
	godotenv.Load()

	return Config{
		Port:          getEnv("PORT", "8080"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "*"),
		DBPath:        getEnv("DB_PATH", "blog.db"),
		AdminToken:    getEnv("ADMIN_TOKEN", ""),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}

}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
