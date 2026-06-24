package dto

import (
	"errors"
	"time"
)

// TelegramAuthRequest is the raw payload produced by the Telegram Login Widget
// and POSTed to /api/v1/auth/telegram. Its Validate method is invoked by the
// BodyValidation middleware before the controller runs.
type TelegramAuthRequest struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

func (r TelegramAuthRequest) Validate() error {
	if r.ID == 0 {
		return errors.New("id is required")
	}
	if r.Hash == "" {
		return errors.New("hash is required")
	}
	if r.AuthDate == 0 {
		return errors.New("auth_date is required")
	}
	return nil
}

// MeResponse is the shape returned by GET /api/v1/auth/me. It extends the
// User schema with an Onboarded flag so the frontend can route new users to
// the onboarding flow without a separate API call.
type MeResponse struct {
	ID         uint      `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	PhotoURL   string    `json:"photo_url"`
	CreatedAt  time.Time `json:"created_at"`
	Onboarded  bool      `json:"onboarded"`
}
