package db

import "time"

// Photo is a single image a user has posted to the public feed. Only the S3
// object key is stored; read URLs are presigned on demand. There are no likes
// or comments by design — the feed is intentionally minimal.
type Photo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Key       string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
