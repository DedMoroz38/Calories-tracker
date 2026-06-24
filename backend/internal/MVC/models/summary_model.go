package models

import (
	"time"

	"calorie-counter/internal/db"

	"gorm.io/gorm"
)

// SummaryModel owns all read-only aggregation queries for the summary and
// stats endpoints. It does not write any data.
type SummaryModel struct {
	DB *gorm.DB
}

// DayTotals holds the aggregated macros for a single calendar day.
type DayTotals struct {
	Date     time.Time
	Calories float64
	Carbs    float64
	Fat      float64
	Protein  float64
}

// TotalsForDay sums all food entries for the user on the given calendar day
// (using the consumed_at column in UTC).
func (m SummaryModel) TotalsForDay(userID uint, day time.Time) (DayTotals, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	var result DayTotals
	result.Date = start

	row := m.DB.
		Model(&db.FoodEntry{}).
		Select("COALESCE(SUM(calories),0) AS calories, COALESCE(SUM(carbs),0) AS carbs, COALESCE(SUM(fat),0) AS fat, COALESCE(SUM(protein),0) AS protein").
		Where("user_id = ? AND consumed_at >= ? AND consumed_at < ?", userID, start, end).
		Row()
	if err := row.Scan(&result.Calories, &result.Carbs, &result.Fat, &result.Protein); err != nil {
		return result, err
	}
	return result, nil
}

// TotalsPerDayInRange returns one DayTotals per calendar day within [start, end),
// ordered ascending. Days with no entries produce a zero-value DayTotals.
func (m SummaryModel) TotalsPerDayInRange(userID uint, start, end time.Time) ([]DayTotals, error) {
	type rawRow struct {
		Day      time.Time
		Calories float64
		Carbs    float64
		Fat      float64
		Protein  float64
	}
	var rows []rawRow
	err := m.DB.
		Model(&db.FoodEntry{}).
		Select("DATE_TRUNC('day', consumed_at) AS day, COALESCE(SUM(calories),0) AS calories, COALESCE(SUM(carbs),0) AS carbs, COALESCE(SUM(fat),0) AS fat, COALESCE(SUM(protein),0) AS protein").
		Where("user_id = ? AND consumed_at >= ? AND consumed_at < ?", userID, start, end).
		Group("day").
		Order("day asc").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// Build a map keyed by truncated day so we can fill missing days with zeros.
	byDay := make(map[time.Time]DayTotals, len(rows))
	for _, r := range rows {
		byDay[r.Day.UTC().Truncate(24*time.Hour)] = DayTotals{
			Date:     r.Day.UTC(),
			Calories: r.Calories,
			Carbs:    r.Carbs,
			Fat:      r.Fat,
			Protein:  r.Protein,
		}
	}

	// Enumerate every calendar day in [start, end) and produce an ordered slice.
	var out []DayTotals
	for d := start; d.Before(end); d = d.Add(24 * time.Hour) {
		key := d.UTC().Truncate(24 * time.Hour)
		if v, ok := byDay[key]; ok {
			out = append(out, v)
		} else {
			out = append(out, DayTotals{Date: d})
		}
	}
	return out, nil
}

// DaysWithEntries returns all distinct calendar days (UTC, truncated) on which
// the user has at least one food entry, ordered descending (most recent first).
// This is used to compute the consecutive-day logging streak.
func (m SummaryModel) DaysWithEntries(userID uint) ([]time.Time, error) {
	var days []time.Time
	err := m.DB.
		Model(&db.FoodEntry{}).
		Select("DISTINCT DATE_TRUNC('day', consumed_at) AS day").
		Where("user_id = ?", userID).
		Order("day desc").
		Pluck("day", &days).Error
	return days, err
}

// WeightEntriesInRange returns weight entries for a user in [start, end),
// ordered ascending by recorded_at.
func (m SummaryModel) WeightEntriesInRange(userID uint, start, end time.Time) ([]db.WeightEntry, error) {
	var entries []db.WeightEntry
	err := m.DB.
		Where("user_id = ? AND recorded_at >= ? AND recorded_at < ?", userID, start, end).
		Order("recorded_at asc").
		Find(&entries).Error
	return entries, err
}
