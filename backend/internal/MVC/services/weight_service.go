package services

import (
	"time"

	"calorie-counter/internal/MVC/models"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"
	"calorie-counter/internal/db"

	"gorm.io/gorm"
)

// WeightService is the use-case layer for the weight-entry domain.
type WeightService struct {
	DB *gorm.DB
}

// Create logs a new weight entry for the user. A missing recorded_at defaults
// to the current time.
func (s WeightService) Create(userID uint, req dto.CreateWeightRequest) (*dto.WeightEntryResponse, *errors.APIError) {
	wm := models.WeightModel{DB: s.DB}

	recordedAt := req.RecordedAt
	if recordedAt.IsZero() {
		recordedAt = time.Now().UTC()
	}

	entry := &db.WeightEntry{
		UserID:     userID,
		Weight:     req.Weight,
		RecordedAt: recordedAt,
	}
	if err := wm.Create(entry); err != nil {
		return nil, errors.Internal("could not log weight entry")
	}

	return &dto.WeightEntryResponse{
		ID:         entry.ID,
		Weight:     entry.Weight,
		RecordedAt: entry.RecordedAt,
	}, nil
}

// List returns all weight entries for the user in chronological order.
func (s WeightService) List(userID uint) ([]dto.WeightEntryResponse, *errors.APIError) {
	wm := models.WeightModel{DB: s.DB}

	entries, err := wm.ListByUser(userID)
	if err != nil {
		return nil, errors.Internal("could not list weight entries")
	}

	out := make([]dto.WeightEntryResponse, len(entries))
	for i, e := range entries {
		out[i] = dto.WeightEntryResponse{
			ID:         e.ID,
			Weight:     e.Weight,
			RecordedAt: e.RecordedAt,
		}
	}
	return out, nil
}
