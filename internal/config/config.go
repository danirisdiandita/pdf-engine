package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DiscordWebhookURL string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		DiscordWebhookURL: getEnv("DISCORD_WEBHOOK_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if err := godotenv.Load("/volume/.env"); err != nil {
		log.Println("Warning: Error loading /volume/.env file")
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
