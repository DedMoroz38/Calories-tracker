package dto

import (
	"errors"
	"time"
)

// CreateWeightRequest is the body for POST /api/v1/weights.
// recorded_at defaults to now when omitted.
type CreateWeightRequest struct {
	Weight     float64   `json:"weight"`
	RecordedAt time.Time `json:"recorded_at"`
}

func (r CreateWeightRequest) Validate() error {
	if r.Weight <= 0 {
		return errors.New("weight must be greater than zero")
	}
	return nil
}

// WeightEntryResponse is the shape of each item in the weights list and the
// response to POST /api/v1/weights.
type WeightEntryResponse struct {
	ID         uint      `json:"id"`
	Weight     float64   `json:"weight"`
	RecordedAt time.Time `json:"recorded_at"`
}
