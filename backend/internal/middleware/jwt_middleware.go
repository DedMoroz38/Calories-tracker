package middleware

import (
	"strings"

	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/config"
	"calorie-counter/internal/util"

	"github.com/gofiber/fiber/v2"
)

// CheckJWT guards protected routes. It accepts the token either as an
// "Authorization: Bearer <token>" header (used by the SPA frontend) or as the
// "jwt_token" cookie, parses it, and stores the claims in c.Locals("jwt_claims")
// for the controller to consume via GetUserID.
func CheckJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := ""
		if header := c.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
			tokenStr = strings.TrimPrefix(header, "Bearer ")
		} else {
			tokenStr = c.Cookies("jwt_token")
		}

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{Message: "missing authentication token"})
		}

		claims, err := util.ParseJWT(tokenStr, config.Values.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{Message: "invalid or expired token"})
		}

		c.Locals("jwt_claims", claims)
		return c.Next()
	}
}

// GetUserID returns the authenticated user id from the JWT claims placed in
// context by CheckJWT, or 0 when the request is not authenticated.
func GetUserID(c *fiber.Ctx) uint {
	claims, ok := c.Locals("jwt_claims").(*util.Claims)
	if !ok {
		return 0
	}
	return claims.UserID
}
