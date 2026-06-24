package controllers

import (
	"calorie-counter/internal/MVC/services"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProfileController struct {
	BaseController
}

// GetProfile handles GET /api/v1/profile. Returns zero-value defaults with
// onboarded=false when the user has not yet completed onboarding.
func (pc ProfileController) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	ps := services.ProfileService{}
	ps.DB = c.Locals("gorm").(*gorm.DB)

	profile, apiErr := ps.GetProfile(userID)
	if apiErr != nil {
		return pc.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: profile})
}

// UpdateProfile handles PUT /api/v1/profile. Used by the onboarding flow and
// subsequent edits. Sets onboarded=true and optionally seeds an initial weight
// entry when current_weight is provided.
func (pc ProfileController) UpdateProfile(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.UpdateProfileRequest)
	userID := middleware.GetUserID(c)

	ps := services.ProfileService{}
	ps.DB = c.Locals("gorm").(*gorm.DB)

	profile, apiErr := ps.UpdateProfile(userID, *req)
	if apiErr != nil {
		return pc.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: profile})
}
