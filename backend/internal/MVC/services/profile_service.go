package services

import (
	stderrors "errors"
	"time"

	"calorie-counter/internal/MVC/models"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"
	"calorie-counter/internal/db"

	"gorm.io/gorm"
)

// ProfileService is the use-case layer for the profile/goals domain.
type ProfileService struct {
	DB *gorm.DB
}

// GetProfile returns the user's profile. When no profile row exists yet (the
// user has not onboarded) a zero-value response with Onboarded=false is
// returned rather than an error.
func (s ProfileService) GetProfile(userID uint) (*dto.ProfileResponse, *errors.APIError) {
	// Identity fields (name + avatar) come from the user row and exist even
	// before onboarding, so resolve them first.
	resp := &dto.ProfileResponse{}
	am := models.AuthModel{DB: s.DB}
	if u, err := am.FindByID(userID); err == nil {
		resp.FirstName = u.FirstName
		resp.Username = u.Username
		resp.AvatarURL = avatarURL(u.AvatarKey, u.PhotoURL)
	}

	pm := models.ProfileModel{DB: s.DB}
	p, err := pm.FindByUserID(userID)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			resp.Onboarded = false
			return resp, nil
		}
		return nil, errors.Internal("could not load profile")
	}
	resp.CalorieGoal = p.CalorieGoal
	resp.CarbsGoal = p.CarbsGoal
	resp.FatGoal = p.FatGoal
	resp.ProteinGoal = p.ProteinGoal
	resp.CurrentWeight = p.CurrentWeight
	resp.GoalWeight = p.GoalWeight
	resp.Direction = p.Direction
	resp.Onboarded = p.Onboarded
	return resp, nil
}

// UpdateProfile upserts the user's profile and sets Onboarded = true. When the
// request includes a non-zero current_weight an initial weight entry is also
// seeded so the weight trend has a starting point.
func (s ProfileService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*dto.ProfileResponse, *errors.APIError) {
	pm := models.ProfileModel{DB: s.DB}

	profile := &db.UserProfile{
		UserID:        userID,
		CalorieGoal:   req.CalorieGoal,
		CarbsGoal:     req.CarbsGoal,
		FatGoal:       req.FatGoal,
		ProteinGoal:   req.ProteinGoal,
		CurrentWeight: req.CurrentWeight,
		GoalWeight:    req.GoalWeight,
		Direction:     req.Direction,
	}
	if err := pm.Upsert(profile); err != nil {
		return nil, errors.Internal("could not save profile")
	}

	// Seed an initial weight entry when the caller supplies current_weight.
	if req.CurrentWeight > 0 {
		wm := models.WeightModel{DB: s.DB}
		entry := &db.WeightEntry{
			UserID:     userID,
			Weight:     req.CurrentWeight,
			RecordedAt: time.Now().UTC(),
		}
		// Non-fatal: if the weight entry fails we still return the saved profile.
		_ = wm.Create(entry)
	}

	return &dto.ProfileResponse{
		CalorieGoal:   profile.CalorieGoal,
		CarbsGoal:     profile.CarbsGoal,
		FatGoal:       profile.FatGoal,
		ProteinGoal:   profile.ProteinGoal,
		CurrentWeight: profile.CurrentWeight,
		GoalWeight:    profile.GoalWeight,
		Direction:     profile.Direction,
		Onboarded:     true,
	}, nil
}

// IsOnboarded returns true when the user has a profile row.
func (s ProfileService) IsOnboarded(userID uint) bool {
	pm := models.ProfileModel{DB: s.DB}
	p, err := pm.FindByUserID(userID)
	if err != nil {
		return false
	}
	return p.Onboarded
}
