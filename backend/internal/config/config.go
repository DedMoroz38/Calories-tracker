package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	TelegramBotToken string
	JWTSecret        string
}

// Values holds the process-wide configuration. It is populated once at startup
// in main() and read by infrastructure that runs outside of an explicit
// dependency chain (DB bootstrap, JWT middleware, auth service). This mirrors
// the "secret set once at startup" pattern described in ARCHITECTURE_FLOW.md.
var Values Config

// LoadConfig reads configuration from the environment (and an optional .env
// file) and also stores the result in the package-level Values global.
func LoadConfig() Config {
	_ = godotenv.Load()
	Values = Config{
		Port:             getEnv("PORT", "3000"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		JWTSecret:        getEnv("JWT_SECRET", "change-me"),
	}
	return Values
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
