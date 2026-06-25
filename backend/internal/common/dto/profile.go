package dto

import "errors"

// UpdateProfileRequest is the body for PUT /api/v1/profile. calorie_goal is
// the only required field; macro goals and weight fields are optional.
type UpdateProfileRequest struct {
	CalorieGoal   float64 `json:"calorie_goal"`
	CarbsGoal     float64 `json:"carbs_goal"`
	FatGoal       float64 `json:"fat_goal"`
	ProteinGoal   float64 `json:"protein_goal"`
	CurrentWeight float64 `json:"current_weight"`
	GoalWeight    float64 `json:"goal_weight"`
	// Direction must be one of: lose | maintain | gain (or empty).
	Direction string `json:"direction"`
}

func (r UpdateProfileRequest) Validate() error {
	if r.CalorieGoal < 0 {
		return errors.New("calorie_goal must be zero or greater")
	}
	switch r.Direction {
	case "", "lose", "maintain", "gain":
		// valid
	default:
		return errors.New("direction must be lose, maintain, or gain")
	}
	return nil
}

// ProfileResponse is the shape returned by GET /api/v1/profile and also
// embedded in the /auth/me response's onboarded check.
type ProfileResponse struct {
	CalorieGoal   float64 `json:"calorie_goal"`
	CarbsGoal     float64 `json:"carbs_goal"`
	FatGoal       float64 `json:"fat_goal"`
	ProteinGoal   float64 `json:"protein_goal"`
	CurrentWeight float64 `json:"current_weight"`
	GoalWeight    float64 `json:"goal_weight"`
	Direction     string  `json:"direction"`
	Onboarded     bool    `json:"onboarded"`
	// Identity fields for the profile page header.
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}
