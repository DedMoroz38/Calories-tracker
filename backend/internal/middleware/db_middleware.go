package middleware

import (
	"calorie-counter/internal/db"

	"github.com/gofiber/fiber/v2"
)

// PGMiddleware stores the single shared *gorm.DB handle in request-local
// context under "gorm". Every controller under /api/v1 reads it back with
// c.Locals("gorm").(*gorm.DB) and passes it down into its service.
func PGMiddleware(dbService *db.DBService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("gorm", dbService.DB)
		return c.Next()
	}
}
