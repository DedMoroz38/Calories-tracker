package routes

import (
	"calorie-counter/internal/db"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterApi is the application's composition root. It creates the shared
// dependencies once (the DB handle), installs the middleware that injects them
// into every request, and registers each domain's route group.
func RegisterApi(app *fiber.App) {
	dbService := db.NewDBService()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Make the shared *gorm.DB available to every /api/v1 request.
	v1.Use(middleware.PGMiddleware(dbService))

	RegisterAuthRoutes(v1.Group("/auth"))
	RegisterFoodRoutes(v1.Group("/foods"))
	RegisterProfileRoutes(v1.Group("/profile"))
	RegisterWeightRoutes(v1.Group("/weights"))
	RegisterSummaryRoutes(v1.Group("/summary"))
	RegisterStatsRoutes(v1.Group("/stats"))
	RegisterPhotoRoutes(v1.Group("/photos"))
	RegisterFeedRoutes(v1.Group("/feed"))
}
