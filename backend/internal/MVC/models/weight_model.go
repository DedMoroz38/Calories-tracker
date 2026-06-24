package models

import (
	"calorie-counter/internal/db"

	"gorm.io/gorm"
)

// WeightModel owns all persistence for the weight-entry domain. The DB handle
// is assigned by the service before any method is called.
type WeightModel struct {
	DB *gorm.DB
}

// Create inserts a new weight entry. On return, entry carries its generated ID.
func (m WeightModel) Create(entry *db.WeightEntry) error {
	return m.DB.Create(entry).Error
}

// ListByUser returns all weight entries for a user in chronological order
// (oldest first) so callers can read the last element as the current weight.
func (m WeightModel) ListByUser(userID uint) ([]db.WeightEntry, error) {
	var entries []db.WeightEntry
	err := m.DB.
		Where("user_id = ?", userID).
		Order("recorded_at asc").
		Find(&entries).Error
	return entries, err
}

// LatestByUser returns the single most-recent weight entry for a user, or
// gorm.ErrRecordNotFound when no entries exist.
func (m WeightModel) LatestByUser(userID uint) (*db.WeightEntry, error) {
	var entry db.WeightEntry
	if err := m.DB.
		Where("user_id = ?", userID).
		Order("recorded_at desc").
		First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}
