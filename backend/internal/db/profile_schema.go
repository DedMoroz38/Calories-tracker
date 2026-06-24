package db

// UserProfile stores a user's goals and physical stats. There is at most one
// row per user (uniqueIndex on UserID). GET /profile returns a default
// zero-value response when no row exists; PUT /profile upserts the row and
// sets Onboarded = true.
type UserProfile struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	UserID        uint    `gorm:"uniqueIndex;not null" json:"user_id"`
	User          User    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CalorieGoal   float64 `json:"calorie_goal"`
	CarbsGoal     float64 `json:"carbs_goal"`
	FatGoal       float64 `json:"fat_goal"`
	ProteinGoal   float64 `json:"protein_goal"`
	CurrentWeight float64 `json:"current_weight"`
	GoalWeight    float64 `json:"goal_weight"`
	// Direction is one of: lose | maintain | gain
	Direction string `json:"direction"`
	Onboarded bool   `json:"onboarded"`
}
