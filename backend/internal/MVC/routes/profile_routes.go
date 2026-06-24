package routes

import (
	"calorie-counter/internal/MVC/controllers"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterProfileRoutes attaches the profile/goals endpoints to a scoped
// router. Every endpoint is protected by JWT auth.
func RegisterProfileRoutes(router fiber.Router) {
	pc := controllers.ProfileController{}

	router.Use(middleware.CheckJWT())

	router.Get("/", pc.GetProfile)
	router.Put("/", middleware.BodyValidation(&dto.UpdateProfileRequest{}), pc.UpdateProfile)
}
