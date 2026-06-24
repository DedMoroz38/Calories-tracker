package dto

// MacroTotals holds calorie and macro sums, used for both daily totals and
// goal values in SummaryResponse.
type MacroTotals struct {
	Calories float64 `json:"calories"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Protein  float64 `json:"protein"`
}

// WeightInfo bundles current and goal weights for the summary endpoint.
type WeightInfo struct {
	Current float64 `json:"current"`
	Goal    float64 `json:"goal"`
}

// DayState is one entry in the 7-day week array returned by GET /summary.
// State is one of: hit | miss | today | future. Fire is true for the current
// day once its calorie goal has been met (drives the streak flame on the dot).
type DayState struct {
	Date   string `json:"date"`
	Letter string `json:"letter"`
	State  string `json:"state"`
	Fire   bool   `json:"fire"`
}

// SummaryResponse is the full payload for GET /api/v1/summary.
type SummaryResponse struct {
	Date   string      `json:"date"`
	Totals MacroTotals `json:"totals"`
	Goals  MacroTotals `json:"goals"`
	Weight WeightInfo  `json:"weight"`
	Streak int         `json:"streak"`
	// StreakActiveToday is true when today's calorie goal has been met, i.e. the
	// streak already counts today (drives the flame next to the date header).
	StreakActiveToday bool       `json:"streak_active_today"`
	Week              []DayState `json:"week"`
}

// CaloriesEntry is one data point in the calories_per_day series.
// Muted is true when no food was logged on that day (value is 0).
type CaloriesEntry struct {
	Label string  `json:"label"`
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Muted bool    `json:"muted"`
}

// WeightTrendEntry is one data point in the weight_trend series.
type WeightTrendEntry struct {
	Label      string `json:"label"`
	Value      float64 `json:"value"`
	RecordedAt string `json:"recorded_at"`
}

// StatsResponse is the full payload for GET /api/v1/stats.
type StatsResponse struct {
	Range          string             `json:"range"`
	CaloriesPerDay []CaloriesEntry    `json:"calories_per_day"`
	WeightTrend    []WeightTrendEntry `json:"weight_trend"`
	AvgPerDay      float64            `json:"avg_per_day"`
	Goal           float64            `json:"goal"`
	DaysUnderGoal  int                `json:"days_under_goal"`
	Streak         int                `json:"streak"`
}
