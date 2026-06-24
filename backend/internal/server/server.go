package server

import (
	"log"

	"calorie-counter/internal/MVC/routes"
	"calorie-counter/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// NewApp builds the Fiber app, installs global middleware, and hands control to
// the composition root. The real application wiring happens in routes.RegisterApi.
func NewApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	routes.RegisterApi(app)
	return app
}

// Serve builds the app and starts listening on the configured port.
func Serve() {
	if err := NewApp().Listen(":" + config.Values.Port); err != nil {
		log.Fatal(err)
	}
}
