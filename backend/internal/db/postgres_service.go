package db

import (
	"log"

	"calorie-counter/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBService holds the single shared *gorm.DB handle. It is created once in the
// composition root (routes.RegisterApi) and injected into every request by
// middleware.PGMiddleware.
type DBService struct {
	DB *gorm.DB
}

// NewDBService opens the PostgreSQL connection from config.Values.DatabaseURL
// and runs AutoMigrate for all schema structs. It fails fast on error because
// the process cannot serve traffic without a database.
func NewDBService() *DBService {
	database, err := gorm.Open(postgres.Open(config.Values.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	if err := database.AutoMigrate(&User{}, &FoodEntry{}, &UserProfile{}, &WeightEntry{}, &Photo{}); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	return &DBService{DB: database}
}
