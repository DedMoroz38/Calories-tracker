package routes

import (
	"calorie-counter/internal/MVC/controllers"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterFoodRoutes attaches the food-entry endpoints to a scoped router. Every
// endpoint here is protected: CheckJWT runs for the whole group.
//
// NOTE: /recent must be registered before /:id so Fiber does not interpret the
// literal string "recent" as a parameter value.
func RegisterFoodRoutes(router fiber.Router) {
	fc := controllers.FoodController{}

	router.Use(middleware.CheckJWT())

	router.Post("/", middleware.BodyValidation(&dto.CreateFoodEntryRequest{}), fc.Create)
	router.Get("/recent", fc.Recent)
	router.Get("/", fc.List)
	router.Delete("/:id", fc.Delete)
}
