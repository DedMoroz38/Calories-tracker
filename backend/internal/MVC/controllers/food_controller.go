package controllers

import (
	"strconv"
	"time"

	"calorie-counter/internal/MVC/services"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FoodController struct {
	BaseController
}

// Create handles POST /api/v1/foods.
func (fc FoodController) Create(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.CreateFoodEntryRequest)
	userID := middleware.GetUserID(c)

	fs := services.FoodService{}
	fs.DB = c.Locals("gorm").(*gorm.DB)

	entry, apiErr := fs.Create(userID, *req)
	if apiErr != nil {
		return fc.Fail(c, apiErr)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{Data: entry})
}

// List handles GET /api/v1/foods with an optional ?date=YYYY-MM-DD filter.
func (fc FoodController) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var day time.Time
	if dateStr := c.Query("date"); dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fc.Fail(c, errors.BadRequest("date must be YYYY-MM-DD"))
		}
		day = t.UTC()
	}

	fs := services.FoodService{}
	fs.DB = c.Locals("gorm").(*gorm.DB)

	entries, apiErr := fs.List(userID, day)
	if apiErr != nil {
		return fc.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: entries})
}

// Recent handles GET /api/v1/foods/recent. It returns the most-recently-used
// distinct named dishes for the user (Quick add and blank names excluded).
func (fc FoodController) Recent(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	fs := services.FoodService{}
	fs.DB = c.Locals("gorm").(*gorm.DB)

	entries, apiErr := fs.Recent(userID)
	if apiErr != nil {
		return fc.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: entries})
}

// Delete handles DELETE /api/v1/foods/:id.
func (fc FoodController) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return fc.Fail(c, errors.BadRequest("invalid food entry id"))
	}

	fs := services.FoodService{}
	fs.DB = c.Locals("gorm").(*gorm.DB)

	if apiErr := fs.Delete(userID, uint(id)); apiErr != nil {
		return fc.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Message: "food entry deleted"})
}
