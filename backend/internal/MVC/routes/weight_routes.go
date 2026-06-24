package routes

import (
	"calorie-counter/internal/MVC/controllers"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterWeightRoutes attaches the weight-entry endpoints to a scoped router.
// Every endpoint is protected by JWT auth.
func RegisterWeightRoutes(router fiber.Router) {
	wc := controllers.WeightController{}

	router.Use(middleware.CheckJWT())

	router.Post("/", middleware.BodyValidation(&dto.CreateWeightRequest{}), wc.Create)
	router.Get("/", wc.List)
}
