package config

import "os"

type Config struct {
	Port          string
	AllowedOrigin string
	DBPath        string
}

func Load() Config {
	cfg := Config{
		Port:          getEnv("PORT", "8080"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "*"),
		DBPath:        getEnv("DB_PATH", "blog.db"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
