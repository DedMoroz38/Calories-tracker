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

	// S3-compatible object storage — used to store user-uploaded photos and
	// avatars. Works with AWS S3 or any S3-compatible service (e.g. a MinIO
	// instance on Railway). When the required vars are unset, the storage layer
	// is disabled and photo endpoints return a 500 telling the operator the
	// bucket is not configured.
	//
	// S3Endpoint is the custom endpoint URL for non-AWS providers. Leave it
	// empty to use real AWS S3; set it (e.g. the public URL of a Railway MinIO
	// service) to target an S3-compatible host. A non-empty endpoint also
	// switches the client to path-style addressing, which MinIO requires.
	S3Endpoint         string
	AWSRegion          string
	AWSS3Bucket        string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
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

		S3Endpoint:         getEnv("S3_ENDPOINT", ""),
		AWSRegion:          getEnv("AWS_REGION", ""),
		AWSS3Bucket:        getEnv("AWS_S3_BUCKET", ""),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
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
