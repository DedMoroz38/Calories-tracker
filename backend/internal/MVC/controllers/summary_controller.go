package controllers

import (
	"calorie-counter/internal/MVC/services"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SummaryController struct {
	BaseController
}

// GetSummary handles GET /api/v1/summary?date=YYYY-MM-DD (date optional, default
// today). Returns totals, goals, weight, streak, and a 7-day dot state array.
func (sc SummaryController) GetSummary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	dateStr := c.Query("date")
	offsetMin := c.QueryInt("tz", 0)

	ss := services.SummaryService{}
	ss.DB = c.Locals("gorm").(*gorm.DB)

	summary, apiErr := ss.GetSummary(userID, dateStr, offsetMin)
	if apiErr != nil {
		return sc.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: summary})
}

// GetStats handles GET /api/v1/stats?range=week|month|year (default week).
// Returns calories-per-day series, weight trend, aggregate stats, and streak.
func (sc SummaryController) GetStats(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	rangeStr := c.Query("range")
	offsetMin := c.QueryInt("tz", 0)

	ss := services.SummaryService{}
	ss.DB = c.Locals("gorm").(*gorm.DB)

	stats, apiErr := ss.GetStats(userID, rangeStr, offsetMin)
	if apiErr != nil {
		return sc.Fail(c, apiErr)
	}

	return c.JSON(dto.BaseResponse{Data: stats})
}
