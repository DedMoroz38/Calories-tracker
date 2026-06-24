package db

import "time"

// User is the schema struct GORM maps to the users table. Each user is keyed by
// their immutable Telegram account id.
type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TelegramID int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	PhotoURL   string    `json:"photo_url"`
	CreatedAt  time.Time `json:"created_at"`
}
