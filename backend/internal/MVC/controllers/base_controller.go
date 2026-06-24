package controllers

import (
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"

	"github.com/gofiber/fiber/v2"
)

// BaseController is embedded by every controller to share HTTP-layer behavior.
// Today that is a single helper for rendering an *APIError as a BaseResponse,
// keeping the error-to-status mapping in one place.
type BaseController struct{}

// Fail writes an APIError using the error's own status code and message.
func (BaseController) Fail(c *fiber.Ctx, err *errors.APIError) error {
	return c.Status(err.StatusCode).JSON(dto.BaseResponse{Message: err.Message})
}
