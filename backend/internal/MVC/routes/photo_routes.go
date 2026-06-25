package routes

import (
	"calorie-counter/internal/MVC/controllers"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterPhotoRoutes attaches the photo-post endpoints. Every endpoint is
// protected by JWT auth.
//
// NOTE: /me must be registered before /:id so Fiber does not treat the literal
// "me" as a parameter value.
func RegisterPhotoRoutes(router fiber.Router) {
	pc := controllers.PhotoController{}

	router.Use(middleware.CheckJWT())

	router.Post("/", pc.Create)
	router.Get("/me", pc.Mine)
	router.Delete("/:id", pc.Delete)
}

// RegisterFeedRoutes attaches the public feed endpoint (other users' photos).
func RegisterFeedRoutes(router fiber.Router) {
	pc := controllers.PhotoController{}

	router.Use(middleware.CheckJWT())

	router.Get("/", pc.Feed)
}
