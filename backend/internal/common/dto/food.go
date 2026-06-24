package dto

import (
	"errors"
	"strings"
	"time"
)

// CreateFoodEntryRequest is the body for POST /api/v1/foods. Its Validate
// method is invoked by the BodyValidation middleware before the controller runs.
type CreateFoodEntryRequest struct {
	Name       string    `json:"name"`
	Calories   float64   `json:"calories"`
	Protein    float64   `json:"protein"`
	Carbs      float64   `json:"carbs"`
	Fat        float64   `json:"fat"`
	ConsumedAt time.Time `json:"consumed_at"`
}

func (r CreateFoodEntryRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	if r.Calories < 0 {
		return errors.New("calories must be zero or greater")
	}
	if r.Protein < 0 || r.Carbs < 0 || r.Fat < 0 {
		return errors.New("macronutrients must be zero or greater")
	}
	return nil
}

// RecentDishResponse is the slim shape returned by GET /api/v1/foods/recent.
// Only the fields needed to re-add a dish from the AddSheet "Recent" tab are
// included; persistence identifiers are intentionally omitted.
type RecentDishResponse struct {
	Name     string  `json:"name"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
}
