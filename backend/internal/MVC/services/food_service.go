package services

import (
	"time"

	"calorie-counter/internal/MVC/models"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"
	"calorie-counter/internal/db"

	"gorm.io/gorm"
)


// FoodService is the use-case layer for food entries. It uses the explicit
// DB-field variant described in ARCHITECTURE_FLOW.md §5; the controller assigns
// the handle with fs.DB = c.Locals("gorm").(*gorm.DB).
type FoodService struct {
	DB *gorm.DB
}

// Create logs a food entry for the given user. A missing ConsumedAt defaults to
// the current time so the client can omit it for "just ate this" logging.
func (s FoodService) Create(userID uint, req dto.CreateFoodEntryRequest) (*db.FoodEntry, *errors.APIError) {
	fm := models.FoodModel{DB: s.DB}

	consumedAt := req.ConsumedAt
	if consumedAt.IsZero() {
		consumedAt = time.Now()
	}

	entry := &db.FoodEntry{
		UserID:     userID,
		Name:       req.Name,
		Calories:   req.Calories,
		Protein:    req.Protein,
		Carbs:      req.Carbs,
		Fat:        req.Fat,
		ConsumedAt: consumedAt,
	}
	if err := fm.Create(entry); err != nil {
		return nil, errors.Internal("could not create food entry")
	}
	return entry, nil
}

// List returns food entries belonging to the user within the half-open window
// [from, to). A zero from or to leaves that bound unconstrained, so passing
// both zero returns every entry.
func (s FoodService) List(userID uint, from, to time.Time) ([]db.FoodEntry, *errors.APIError) {
	fm := models.FoodModel{DB: s.DB}

	entries, err := fm.ListByUser(userID, from, to)
	if err != nil {
		return nil, errors.Internal("could not list food entries")
	}
	return entries, nil
}

// Recent returns distinct recent named dishes (excluding blank names and
// "Quick add") for the user, capped at 20, most-recently-consumed first.
// It projects to the slim RecentDishResponse shape — persistence identifiers
// are excluded so the client only sees what it needs to re-add a dish.
func (s FoodService) Recent(userID uint) ([]dto.RecentDishResponse, *errors.APIError) {
	fm := models.FoodModel{DB: s.DB}

	entries, err := fm.ListRecent(userID)
	if err != nil {
		return nil, errors.Internal("could not list recent dishes")
	}

	out := make([]dto.RecentDishResponse, len(entries))
	for i, e := range entries {
		out[i] = dto.RecentDishResponse{
			Name:     e.Name,
			Calories: e.Calories,
			Protein:  e.Protein,
			Carbs:    e.Carbs,
			Fat:      e.Fat,
		}
	}
	return out, nil
}

// Delete removes a single food entry owned by the user. It reports NotFound when
// no matching row exists so a user cannot probe other users' entries.
func (s FoodService) Delete(userID, id uint) *errors.APIError {
	fm := models.FoodModel{DB: s.DB}

	affected, err := fm.DeleteByIDAndUser(id, userID)
	if err != nil {
		return errors.Internal("could not delete food entry")
	}
	if affected == 0 {
		return errors.NotFound("food entry not found")
	}
	return nil
}
