package services

import (
	"fmt"
	"time"

	"calorie-counter/internal/MVC/models"
	"calorie-counter/internal/common/dto"
	"calorie-counter/internal/common/errors"

	"gorm.io/gorm"
)

// SummaryService is the use-case layer for summary and stats endpoints.
type SummaryService struct {
	DB *gorm.DB
}

// dayLetters maps time.Weekday to a single uppercase letter.
var dayLetters = [7]string{"S", "M", "T", "W", "T", "F", "S"}

// computeStreak counts consecutive calendar days going backwards from today
// on which the user logged at least one food entry.
func computeStreak(days []time.Time) int {
	if len(days) == 0 {
		return 0
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	streak := 0
	expected := today
	for _, d := range days {
		day := d.UTC().Truncate(24 * time.Hour)
		if day.Equal(expected) {
			streak++
			expected = expected.Add(-24 * time.Hour)
		} else if day.Before(expected) {
			break
		}
	}
	return streak
}

// GetSummary builds the home-screen payload for the given date (default today).
func (s SummaryService) GetSummary(userID uint, dateStr string) (*dto.SummaryResponse, *errors.APIError) {
	sm := models.SummaryModel{DB: s.DB}
	pm := models.ProfileModel{DB: s.DB}
	wm := models.WeightModel{DB: s.DB}

	// Resolve the target date.
	var target time.Time
	if dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, errors.BadRequest("date must be YYYY-MM-DD")
		}
		target = t.UTC()
	} else {
		target = time.Now().UTC()
	}
	targetDate := target.Truncate(24 * time.Hour)

	// Day totals for the target date.
	totals, err := sm.TotalsForDay(userID, targetDate)
	if err != nil {
		return nil, errors.Internal("could not compute daily totals")
	}

	// User goals from profile (zero if not yet set).
	var calorieGoal, carbsGoal, fatGoal, proteinGoal float64
	var goalWeight float64
	if p, err := pm.FindByUserID(userID); err == nil {
		calorieGoal = p.CalorieGoal
		carbsGoal = p.CarbsGoal
		fatGoal = p.FatGoal
		proteinGoal = p.ProteinGoal
		goalWeight = p.GoalWeight
	}

	// Current weight from latest weight entry.
	var currentWeight float64
	if latest, err := wm.LatestByUser(userID); err == nil {
		currentWeight = latest.Weight
	}

	// Streak: consecutive days with any logged food, counting back from today.
	activeDays, err := sm.DaysWithEntries(userID)
	if err != nil {
		return nil, errors.Internal("could not compute streak")
	}
	streak := computeStreak(activeDays)

	// Build a set of days that have food entries for O(1) lookup below.
	loggedDays := make(map[string]bool, len(activeDays))
	for _, d := range activeDays {
		loggedDays[d.UTC().Truncate(24*time.Hour).Format("2006-01-02")] = true
	}

	// 7-day week ending on targetDate.
	today := time.Now().UTC().Truncate(24 * time.Hour)
	week := make([]dto.DayState, 7)
	for i := range 7 {
		d := targetDate.Add(time.Duration(i-6) * 24 * time.Hour)
		key := d.Format("2006-01-02")
		letter := dayLetters[d.Weekday()]
		var state string
		switch {
		case d.After(today):
			state = "future"
		case d.Equal(today):
			state = "today"
		case loggedDays[key]:
			state = "hit"
		default:
			state = "miss"
		}
		week[i] = dto.DayState{Date: key, Letter: letter, State: state}
	}

	return &dto.SummaryResponse{
		Date: targetDate.Format("2006-01-02"),
		Totals: dto.MacroTotals{
			Calories: totals.Calories,
			Carbs:    totals.Carbs,
			Fat:      totals.Fat,
			Protein:  totals.Protein,
		},
		Goals: dto.MacroTotals{
			Calories: calorieGoal,
			Carbs:    carbsGoal,
			Fat:      fatGoal,
			Protein:  proteinGoal,
		},
		Weight: dto.WeightInfo{
			Current: currentWeight,
			Goal:    goalWeight,
		},
		Streak: streak,
		Week:   week,
	}, nil
}

// GetStats builds the stats-screen payload for week / month / year ranges.
//
// week  → 7 daily data points
// month → 30 daily data points
// year  → 12 monthly aggregated data points
func (s SummaryService) GetStats(userID uint, rangeStr string) (*dto.StatsResponse, *errors.APIError) {
	switch rangeStr {
	case "", "week", "month", "year":
	default:
		return nil, errors.BadRequest("range must be week, month, or year")
	}
	if rangeStr == "" {
		rangeStr = "week"
	}

	sm := models.SummaryModel{DB: s.DB}
	pm := models.ProfileModel{DB: s.DB}

	var calorieGoal float64
	if p, err := pm.FindByUserID(userID); err == nil {
		calorieGoal = p.CalorieGoal
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)
	var start, end time.Time

	switch rangeStr {
	case "week":
		start = now.Add(-6 * 24 * time.Hour)
		end = now.Add(24 * time.Hour)
	case "month":
		start = now.Add(-29 * 24 * time.Hour)
		end = now.Add(24 * time.Hour)
	case "year":
		// 12 months ending with (and including) the current month.
		// Subtract 11 months from the start of this month to get the first month.
		start = time.Date(now.Year(), now.Month()-11, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	}

	// Streak (calculated before the range-specific logic).
	activeDays, err := sm.DaysWithEntries(userID)
	if err != nil {
		return nil, errors.Internal("could not compute streak")
	}
	streak := computeStreak(activeDays)

	if rangeStr == "year" {
		return s.buildYearStats(sm, userID, start, end, calorieGoal, streak)
	}

	// week / month: daily granularity.
	dayTotals, err := sm.TotalsPerDayInRange(userID, start, end)
	if err != nil {
		return nil, errors.Internal("could not compute calories series")
	}

	caloriesPerDay := make([]dto.CaloriesEntry, len(dayTotals))
	var total float64
	daysUnderGoal := 0
	for i, t := range dayTotals {
		label := labelForDay(rangeStr, t.Date, i)
		caloriesPerDay[i] = dto.CaloriesEntry{
			Label: label,
			Date:  t.Date.Format("2006-01-02"),
			Value: t.Calories,
			Muted: t.Calories == 0,
		}
		total += t.Calories
		if t.Calories > 0 && t.Calories < calorieGoal {
			daysUnderGoal++
		}
	}

	var avgPerDay float64
	if len(dayTotals) > 0 {
		avgPerDay = total / float64(len(dayTotals))
	}

	// Weight trend: individual entries in range.
	weightEntries, err := sm.WeightEntriesInRange(userID, start, end)
	if err != nil {
		return nil, errors.Internal("could not compute weight trend")
	}
	weightTrend := make([]dto.WeightTrendEntry, len(weightEntries))
	for i, we := range weightEntries {
		weightTrend[i] = dto.WeightTrendEntry{
			Label:      labelForDay(rangeStr, we.RecordedAt.UTC(), i),
			Value:      we.Weight,
			RecordedAt: we.RecordedAt.UTC().Format(time.RFC3339),
		}
	}

	return &dto.StatsResponse{
		Range:          rangeStr,
		CaloriesPerDay: caloriesPerDay,
		WeightTrend:    weightTrend,
		AvgPerDay:      avgPerDay,
		Goal:           calorieGoal,
		DaysUnderGoal:  daysUnderGoal,
		Streak:         streak,
	}, nil
}

// buildYearStats produces 12 monthly-aggregated data points for the year range.
func (s SummaryService) buildYearStats(
	sm models.SummaryModel,
	userID uint,
	start, end time.Time,
	calorieGoal float64,
	streak int,
) (*dto.StatsResponse, *errors.APIError) {
	// Fetch raw daily totals and aggregate into months ourselves so we don't
	// need a separate SQL query.
	dayTotals, err := sm.TotalsPerDayInRange(userID, start, end)
	if err != nil {
		return nil, errors.Internal("could not compute yearly calories")
	}

	// Build a map: year-month → aggregated calories & count.
	type monthKey struct{ year int; month time.Month }
	type monthAgg struct {
		calories float64
		days     int
		date     time.Time
	}
	monthMap := make(map[monthKey]*monthAgg)
	var monthOrder []monthKey

	// Enumerate 12 months in order so we always emit all 12 even with no data.
	for i := range 12 {
		t := time.Date(start.Year(), start.Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		k := monthKey{t.Year(), t.Month()}
		if _, ok := monthMap[k]; !ok {
			monthMap[k] = &monthAgg{date: t}
			monthOrder = append(monthOrder, k)
		}
	}
	for _, d := range dayTotals {
		k := monthKey{d.Date.Year(), d.Date.Month()}
		if agg, ok := monthMap[k]; ok {
			agg.calories += d.Calories
			if d.Calories > 0 {
				agg.days++
			}
		}
	}

	caloriesPerDay := make([]dto.CaloriesEntry, len(monthOrder))
	var totalCal float64
	daysUnderGoal := 0
	for i, k := range monthOrder {
		agg := monthMap[k]
		caloriesPerDay[i] = dto.CaloriesEntry{
			Label: agg.date.Format("Jan"),
			Date:  agg.date.Format("2006-01-02"),
			Value: agg.calories,
			Muted: agg.calories == 0,
		}
		totalCal += agg.calories
		if agg.calories > 0 && agg.calories < calorieGoal*float64(agg.days) {
			daysUnderGoal++
		}
	}

	var avgPerDay float64
	activeDayCnt := 0
	for _, d := range dayTotals {
		if d.Calories > 0 {
			activeDayCnt++
			avgPerDay += d.Calories
		}
	}
	if activeDayCnt > 0 {
		avgPerDay /= float64(activeDayCnt)
	}

	// Weight trend: one entry per month (average or latest weight in that month).
	weightEntries, err := sm.WeightEntriesInRange(userID, start, end)
	if err != nil {
		return nil, errors.Internal("could not compute weight trend")
	}

	// Keep only the last weight entry per month.
	latestWeight := make(map[monthKey]dto.WeightTrendEntry)
	for _, we := range weightEntries {
		k := monthKey{we.RecordedAt.UTC().Year(), we.RecordedAt.UTC().Month()}
		latestWeight[k] = dto.WeightTrendEntry{
			Label:      we.RecordedAt.UTC().Format("Jan"),
			Value:      we.Weight,
			RecordedAt: we.RecordedAt.UTC().Format(time.RFC3339),
		}
	}
	var weightTrend []dto.WeightTrendEntry
	for _, k := range monthOrder {
		if entry, ok := latestWeight[k]; ok {
			weightTrend = append(weightTrend, entry)
		}
	}

	return &dto.StatsResponse{
		Range:          "year",
		CaloriesPerDay: caloriesPerDay,
		WeightTrend:    weightTrend,
		AvgPerDay:      avgPerDay,
		Goal:           calorieGoal,
		DaysUnderGoal:  daysUnderGoal,
		Streak:         streak,
	}, nil
}

// labelForDay produces a short display label for a data point depending on
// the range type: single letter for week, day-of-month for month, month
// abbreviation for year.
func labelForDay(rangeStr string, d time.Time, _ int) string {
	switch rangeStr {
	case "week":
		return dayLetters[d.Weekday()]
	case "month":
		return fmt.Sprintf("%d", d.Day())
	default: // year
		return d.Format("Jan")
	}
}
