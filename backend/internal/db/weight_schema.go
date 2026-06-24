package db

import "time"

// WeightEntry is the schema struct GORM maps to the weight_entries table. Each
// row is a single logged body-weight measurement belonging to a user. The most
// recent entry is used as the user's current weight in summaries.
type WeightEntry struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Weight     float64   `json:"weight"`
	RecordedAt time.Time `gorm:"index" json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
}
