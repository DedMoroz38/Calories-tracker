package services

import (
	"testing"
	"time"

	"calorie-counter/internal/MVC/models"
)

// UTC+3 client: getTimezoneOffset() returns -180.
const tzPlus3 = -180

func entry(localWall string, cals float64, offsetMin int) models.EntryCalories {
	// localWall is "2006-01-02 15:04" in the client's local wall clock; convert
	// to the absolute UTC instant the DB would store.
	t, err := time.Parse("2006-01-02 15:04", localWall)
	if err != nil {
		panic(err)
	}
	utc := t.Add(time.Duration(offsetMin) * time.Minute) // local -> UTC
	return models.EntryCalories{ConsumedAt: utc, Calories: cals}
}

func TestLocalStreakDay_3amRollover(t *testing.T) {
	day := func(localWall string) string {
		return localStreakDay(entry(localWall, 0, tzPlus3).ConsumedAt, tzPlus3).Format("2006-01-02")
	}
	cases := map[string]string{
		"2026-06-25 12:00": "2026-06-25", // midday -> same day
		"2026-06-25 02:59": "2026-06-24", // before 3am -> previous day
		"2026-06-25 03:00": "2026-06-25", // 3am sharp -> new day
		"2026-06-26 01:30": "2026-06-25", // 1:30am -> previous day
		"2026-06-25 23:59": "2026-06-25",
	}
	for wall, want := range cases {
		if got := day(wall); got != want {
			t.Errorf("local wall %s: got streak-day %s, want %s", wall, got, want)
		}
	}
}

// Values read back from Postgres can carry a non-UTC location. localStreakDay
// must normalise to UTC so the day bucket doesn't shift by the server's zone.
func TestLocalStreakDay_NonUTCInput(t *testing.T) {
	utc := time.Date(2026, 6, 23, 22, 15, 0, 0, time.UTC)
	moscow := time.FixedZone("MSK", 3*3600)
	want := localStreakDay(utc, tzPlus3)
	if got := localStreakDay(utc.In(moscow), tzPlus3); !got.Equal(want) {
		t.Errorf("non-UTC input gave %s, want %s", got.Format("2006-01-02"), want.Format("2006-01-02"))
	}
}

func TestComputeCalorieStreak(t *testing.T) {
	const goal = 1800.0
	// "now" = 2026-06-25 14:00 local (UTC+3).
	now := entry("2026-06-25 14:00", 0, tzPlus3).ConsumedAt
	currentDay := localStreakDay(now, tzPlus3)

	tests := []struct {
		name       string
		entries    []models.EntryCalories
		goal       float64
		wantStreak int
		wantToday  bool
	}{
		{
			name: "today met + yesterday met, 2 days ago missed -> streak 2, fire today",
			entries: []models.EntryCalories{
				entry("2026-06-25 13:00", 2000, tzPlus3), // today, met
				entry("2026-06-24 20:00", 1900, tzPlus3), // yesterday, met
				entry("2026-06-23 20:00", 500, tzPlus3),  // 2 days ago, missed
			},
			goal: goal, wantStreak: 2, wantToday: true,
		},
		{
			name: "today not yet met but yesterday+before met -> streak 2, no fire",
			entries: []models.EntryCalories{
				entry("2026-06-25 13:00", 500, tzPlus3),  // today, not met
				entry("2026-06-24 20:00", 1900, tzPlus3), // yesterday, met
				entry("2026-06-23 20:00", 1850, tzPlus3), // 2 days ago, met
				entry("2026-06-22 20:00", 100, tzPlus3),  // missed
			},
			goal: goal, wantStreak: 2, wantToday: false,
		},
		{
			name: "late-night entry before 3am counts toward previous day",
			entries: []models.EntryCalories{
				entry("2026-06-25 13:00", 2000, tzPlus3), // today met
				entry("2026-06-25 01:00", 1000, tzPlus3), // 1am -> belongs to Jun 24
				entry("2026-06-24 22:00", 1000, tzPlus3), // Jun 24, +1000 above = 2000 met
			},
			goal: goal, wantStreak: 2, wantToday: true,
		},
		{
			name:       "no goal set -> no streak",
			entries:    []models.EntryCalories{entry("2026-06-25 13:00", 5000, tzPlus3)},
			goal:       0,
			wantStreak: 0, wantToday: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bucket := bucketCaloriesByStreakDay(tc.entries, tzPlus3)
			streak, today, days := computeCalorieStreak(bucket, tc.goal, currentDay)
			if streak != tc.wantStreak {
				t.Errorf("streak = %d, want %d", streak, tc.wantStreak)
			}
			if len(days) != tc.wantStreak {
				t.Errorf("streakDays count = %d, want %d", len(days), tc.wantStreak)
			}
			if today != tc.wantToday {
				t.Errorf("todayMet = %v, want %v", today, tc.wantToday)
			}
		})
	}
}
