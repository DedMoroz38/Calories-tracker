package middleware

import (
	"reflect"

	"calorie-counter/internal/common/dto"

	"github.com/gofiber/fiber/v2"
)

// Validatable is implemented by every request DTO in internal/common/dto.
type Validatable interface {
	Validate() error
}

// BodyValidation parses the JSON request body into a fresh instance of the
// supplied DTO prototype, runs its Validate method, and stores the result in
// c.Locals("validatedBody"). The controller then reads it back with a type
// assertion. A new instance is allocated per request via reflection so the
// prototype is never shared across concurrent requests.
func BodyValidation(prototype Validatable) fiber.Handler {
	t := reflect.TypeOf(prototype)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return func(c *fiber.Ctx) error {
		req := reflect.New(t).Interface().(Validatable)

		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{Message: "invalid request body"})
		}

		if err := req.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.BaseResponse{Message: err.Error()})
		}

		c.Locals("validatedBody", req)
		return c.Next()
	}
}
