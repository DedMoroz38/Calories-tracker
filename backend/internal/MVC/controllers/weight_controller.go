package controllers

import (
	"calorie-counter/internal/MVC/services"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WeightController struct {
	BaseController
}

// Create handles POST /api/v1/weights.
func (wc WeightController) Create(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.CreateWeightRequest)
	userID := middleware.GetUserID(c)

	ws := services.WeightService{}
	ws.DB = c.Locals("gorm").(*gorm.DB)

	entry, apiErr := ws.Create(userID, *req)
	if apiErr != nil {
		return wc.Fail(c, apiErr)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{Data: entry})
}

// List handles GET /api/v1/weights. Returns all weight entries for the user in
// chronological order (oldest first); the last element is the current weight.
func (wc WeightController) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	ws := services.WeightService{}
	ws.DB = c.Locals("gorm").(*gorm.DB)

	entries, apiErr := ws.List(userID)
	if apiErr != nil {
		return wc.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: entries})
}
