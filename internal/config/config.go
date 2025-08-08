package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token          string
	WeatherAPIKey  string
	WeatherBaseURL string
	LogLevel       string
}

func Load() (Config, error) {
	// Try load .env if present
	_ = godotenv.Load(".env")

	cfg := Config{
		Token:          os.Getenv("TOKEN"),
		WeatherAPIKey:  os.Getenv("WEATHER_API_KEY"),
		WeatherBaseURL: getenvDefault("WEATHER_BASE_URL", "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"),
		LogLevel:       getenvDefault("LOG_LEVEL", "info"),
	}

	if cfg.Token == "" {
		return Config{}, fmt.Errorf("TOKEN is required")
	}
	if cfg.WeatherAPIKey == "" {
		return Config{}, fmt.Errorf("WEATHER_API_KEY is required")
	}

	return cfg, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
