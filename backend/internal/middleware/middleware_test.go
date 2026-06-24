package middleware

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/config"
	"calorie-counter/internal/util"

	"github.com/gofiber/fiber/v2"
)

func postJSON(t *testing.T, app *fiber.App, target, body string) (int, string) {
	t.Helper()
	req := httptest.NewRequest("POST", target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	out, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(out)
}

func TestBodyValidation(t *testing.T) {
	app := fiber.New()
	app.Post("/foods", BodyValidation(&dto.CreateFoodEntryRequest{}), func(c *fiber.Ctx) error {
		req := c.Locals("validatedBody").(*dto.CreateFoodEntryRequest)
		return c.SendString(req.Name)
	})

	// Valid body reaches the handler and is available in locals.
	if status, body := postJSON(t, app, "/foods", `{"name":"Apple","calories":95}`); status != fiber.StatusOK || body != "Apple" {
		t.Fatalf("valid body: status=%d body=%q, want 200/Apple", status, body)
	}

	// Missing required name is rejected by Validate before the handler runs.
	if status, _ := postJSON(t, app, "/foods", `{"calories":95}`); status != fiber.StatusBadRequest {
		t.Fatalf("invalid body: status=%d, want 400", status)
	}
}

func TestCheckJWT(t *testing.T) {
	config.Values.JWTSecret = "unit-secret"

	app := fiber.New()
	app.Get("/me", CheckJWT(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"user_id": GetUserID(c)})
	})

	// No token -> 401.
	noTok, _ := app.Test(httptest.NewRequest("GET", "/me", nil))
	if noTok.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("no token: status = %d, want 401", noTok.StatusCode)
	}

	// Valid Bearer token -> 200 and GetUserID returns the claim.
	token, _ := util.GenerateJWT(7, config.Values.JWTSecret)
	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	ok, _ := app.Test(req)
	if ok.StatusCode != fiber.StatusOK {
		t.Fatalf("valid token: status = %d, want 200", ok.StatusCode)
	}
	body, _ := io.ReadAll(ok.Body)
	if !strings.Contains(string(body), `"user_id":7`) {
		t.Fatalf("body = %q, want user_id 7", string(body))
	}
}
