package routes

import (
	"calorie-counter/internal/MVC/controllers"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterAuthRoutes attaches the auth endpoints to a scoped router. The router
// itself (and the /api/v1 prefix) is created by the composition root.
func RegisterAuthRoutes(router fiber.Router) {
	ac := controllers.AuthController{}

	router.Post("/telegram", middleware.BodyValidation(&dto.TelegramAuthRequest{}), ac.TelegramLogin)
	router.Get("/me", middleware.CheckJWT(), ac.Me)
}
