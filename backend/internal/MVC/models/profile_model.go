package models

import (
	"calorie-counter/internal/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProfileModel owns all persistence for the user-profile/goals domain. The DB
// handle is assigned by the service before any method is called.
type ProfileModel struct {
	DB *gorm.DB
}

// FindByUserID loads the user's profile row. Returns gorm.ErrRecordNotFound
// when the user has not yet completed onboarding.
func (m ProfileModel) FindByUserID(userID uint) (*db.UserProfile, error) {
	var p db.UserProfile
	if err := m.DB.Where("user_id = ?", userID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// Upsert creates or fully replaces the profile row for a user and sets
// Onboarded = true. It uses PostgreSQL's ON CONFLICT (user_id) DO UPDATE so
// a single statement handles both the first-time and subsequent updates.
func (m ProfileModel) Upsert(p *db.UserProfile) error {
	p.Onboarded = true
	return m.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"calorie_goal", "carbs_goal", "fat_goal", "protein_goal",
				"current_weight", "goal_weight", "direction", "onboarded",
			}),
		}).
		Create(p).Error
}
