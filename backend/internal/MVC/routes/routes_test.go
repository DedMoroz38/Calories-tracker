package routes

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// newTestApp wires the real route groups under /api/v1 but injects an empty
// *gorm.DB into request context in place of middleware.PGMiddleware, so the
// controller's c.Locals("gorm").(*gorm.DB) assertion succeeds. The paths
// exercised here (validation + auth failures) return before touching the DB.
func newTestApp() *fiber.App {
	app := fiber.New()
	v1 := app.Group("/api/v1")
	v1.Use(func(c *fiber.Ctx) error {
		c.Locals("gorm", &gorm.DB{})
		return c.Next()
	})
	RegisterAuthRoutes(v1.Group("/auth"))
	RegisterFoodRoutes(v1.Group("/foods"))
	return app
}

func TestRoutesRegistered(t *testing.T) {
	app := newTestApp()

	want := map[string]bool{
		"POST /api/v1/auth/telegram": false,
		"GET /api/v1/auth/me":        false,
		"POST /api/v1/foods":         false,
		"GET /api/v1/foods":          false,
		"DELETE /api/v1/foods/:id":   false,
	}
	for _, r := range app.GetRoutes() {
		key := r.Method + " " + r.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for route, found := range want {
		if !found {
			t.Errorf("route not registered: %s", route)
		}
	}
}

func TestTelegramLoginRejectsBadHash(t *testing.T) {
	app := newTestApp()

	// Well-formed body (passes BodyValidation) but the hash is invalid, so the
	// service returns Unauthorized before any DB access.
	req := httptest.NewRequest("POST", "/api/v1/auth/telegram",
		strings.NewReader(`{"id":1,"first_name":"X","auth_date":9999999999,"hash":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 401 (body: %s)", resp.StatusCode, body)
	}
}

func TestTelegramLoginRejectsMissingFields(t *testing.T) {
	app := newTestApp()

	req := httptest.NewRequest("POST", "/api/v1/auth/telegram",
		strings.NewReader(`{"first_name":"X"}`)) // no id/hash/auth_date
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestFoodRoutesRequireAuth(t *testing.T) {
	app := newTestApp()

	for _, tc := range []struct{ method, path string }{
		{"GET", "/api/v1/foods"},
		{"POST", "/api/v1/foods"},
		{"DELETE", "/api/v1/foods/1"},
	} {
		resp, _ := app.Test(httptest.NewRequest(tc.method, tc.path, nil))
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("%s %s: status = %d, want 401", tc.method, tc.path, resp.StatusCode)
		}
	}
}
