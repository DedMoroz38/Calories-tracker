package models

import (
	"time"

	"calorie-counter/internal/db"

	"gorm.io/gorm"
)

// FoodModel owns all persistence for the food-entry domain. The DB handle is
// assigned by the service (fm.DB = fs.DB) before any method is called.
type FoodModel struct {
	DB *gorm.DB
}

// Create inserts a new food entry. On return, entry carries its generated ID.
func (m FoodModel) Create(entry *db.FoodEntry) error {
	return m.DB.Create(entry).Error
}

// ListByUser returns a user's food entries, most recently consumed first. When
// day is non-zero only entries whose consumed_at falls on that calendar day
// (UTC) are returned.
func (m FoodModel) ListByUser(userID uint, day time.Time) ([]db.FoodEntry, error) {
	var entries []db.FoodEntry
	q := m.DB.Where("user_id = ?", userID)
	if !day.IsZero() {
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
		end := start.Add(24 * time.Hour)
		q = q.Where("consumed_at >= ? AND consumed_at < ?", start, end)
	}
	err := q.Order("consumed_at desc").Find(&entries).Error
	return entries, err
}

// ListRecent returns the single most-recent entry for each distinct dish name,
// excluding blank names and the literal string "Quick add", limited to 20
// results ordered by most-recently-consumed. It uses a PostgreSQL
// DISTINCT ON expression and is intentionally PostgreSQL-specific.
func (m FoodModel) ListRecent(userID uint) ([]db.FoodEntry, error) {
	var entries []db.FoodEntry
	err := m.DB.Raw(`
		WITH latest AS (
			SELECT DISTINCT ON (name) id, user_id, name, calories, protein, carbs, fat, consumed_at, created_at
			FROM food_entries
			WHERE user_id = ? AND name != '' AND name != 'Quick add'
			ORDER BY name, consumed_at DESC
		)
		SELECT * FROM latest ORDER BY consumed_at DESC LIMIT 20
	`, userID).Scan(&entries).Error
	return entries, err
}

// DeleteByIDAndUser removes a single entry scoped to its owner and reports how
// many rows were affected (0 means the entry did not exist for this user).
func (m FoodModel) DeleteByIDAndUser(id, userID uint) (int64, error) {
	res := m.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&db.FoodEntry{})
	return res.RowsAffected, res.Error
}
