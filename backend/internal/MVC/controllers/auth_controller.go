package controllers

import (
	"calorie-counter/internal/MVC/services"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct {
	BaseController
}

// TelegramLogin handles POST /api/v1/auth/telegram.
func (ac AuthController) TelegramLogin(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.TelegramAuthRequest)

	as := services.AuthService{}
	as.DB = c.Locals("gorm").(*gorm.DB)

	token, apiErr := as.TelegramLogin(*req)
	if apiErr != nil {
		return ac.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: fiber.Map{"token": token}})
}

// Me handles GET /api/v1/auth/me and returns the authenticated user's profile
// extended with an onboarded flag.
func (ac AuthController) Me(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	as := services.AuthService{}
	as.DB = c.Locals("gorm").(*gorm.DB)

	user, apiErr := as.GetUser(userID)
	if apiErr != nil {
		return ac.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: user})
}
