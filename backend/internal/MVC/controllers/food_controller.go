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

// List handles GET /api/v1/foods. The client supplies the half-open window it
// cares about as absolute instants — ?from=<RFC3339>&to=<RFC3339> — typically
// the user's local calendar day. Both are optional; when omitted that side of
// the window is unbounded, so no params returns the user's full history. The
// server makes no timezone assumptions: it only compares against consumed_at.
func (fc FoodController) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var from, to time.Time
	if s := c.Query("from"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return fc.Fail(c, errors.BadRequest("from must be an RFC3339 timestamp"))
		}
		from = t
	}
	if s := c.Query("to"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return fc.Fail(c, errors.BadRequest("to must be an RFC3339 timestamp"))
		}
		to = t
	}

	fs := services.FoodService{}
	fs.DB = c.Locals("gorm").(*gorm.DB)

	entries, apiErr := fs.List(userID, from, to)
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
