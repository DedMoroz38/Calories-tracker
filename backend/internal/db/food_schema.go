package db

import "time"

// FoodEntry is the schema struct GORM maps to the food_entries table. Each row
// is a single logged item of food belonging to a user.
type FoodEntry struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name       string    `gorm:"not null" json:"name"`
	Calories   float64   `json:"calories"`
	Protein    float64   `json:"protein"`
	Carbs      float64   `json:"carbs"`
	Fat        float64   `json:"fat"`
	ConsumedAt time.Time `gorm:"index" json:"consumed_at"`
	CreatedAt  time.Time `json:"created_at"`
}
