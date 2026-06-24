package routes

import (
	"calorie-counter/internal/MVC/controllers"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterSummaryRoutes attaches the summary and stats endpoints to a scoped
// router. Both endpoints are protected by JWT auth.
func RegisterSummaryRoutes(router fiber.Router) {
	sc := controllers.SummaryController{}

	router.Use(middleware.CheckJWT())

	router.Get("/", sc.GetSummary)
}

// RegisterStatsRoutes attaches the stats endpoint to a scoped router.
func RegisterStatsRoutes(router fiber.Router) {
	sc := controllers.SummaryController{}

	router.Use(middleware.CheckJWT())

	router.Get("/", sc.GetStats)
}
