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

// streakRolloverHours is how far past local midnight a "streak day" stays open.
// Food logged before 3am counts toward the previous day, and today's day does
// not roll over until 3am — so late-night eating doesn't break a streak.
const streakRolloverHours = 3

// localStreakDay maps an absolute instant to the date-key of the local,
// 3am-anchored "streak day" it belongs to. offsetMin is the client's
// JavaScript getTimezoneOffset() value (minutes to add to local time to reach
// UTC; e.g. UTC+3 → -180). The returned key is a midnight-UTC time used only
// as a stable map/identity key — it carries no timezone meaning itself.
func localStreakDay(instant time.Time, offsetMin int) time.Time {
	// Normalise to UTC first: values read back from the DB may carry the
	// server's local zone, and the wall-clock components below must be derived
	// from a known reference, not whatever zone the time happens to be in.
	utc := instant.UTC()
	local := utc.Add(time.Duration(-offsetMin) * time.Minute)
	shifted := local.Add(-streakRolloverHours * time.Hour)
	return time.Date(shifted.Year(), shifted.Month(), shifted.Day(), 0, 0, 0, 0, time.UTC)
}

// bucketCaloriesByStreakDay sums each entry's calories into its local
// 3am-anchored streak day.
func bucketCaloriesByStreakDay(entries []models.EntryCalories, offsetMin int) map[time.Time]float64 {
	bucket := make(map[time.Time]float64, len(entries))
	for _, e := range entries {
		bucket[localStreakDay(e.ConsumedAt, offsetMin)] += e.Calories
	}
	return bucket
}

// computeCalorieStreak counts consecutive streak days, ending at currentDay,
// whose total calories reached the goal. todayMet reports whether currentDay
// itself has already hit the goal. An unmet current day does not break the
// streak (the day isn't over yet) — it simply isn't counted, and the run of
// completed days before it is still returned.
func computeCalorieStreak(bucket map[time.Time]float64, goal float64, currentDay time.Time) (streak int, todayMet bool, streakDays map[time.Time]bool) {
	streakDays = map[time.Time]bool{}
	if goal <= 0 {
		return 0, false, streakDays
	}
	todayMet = bucket[currentDay] >= goal

	day := currentDay
	if !todayMet {
		day = day.AddDate(0, 0, -1)
	}
	for bucket[day] >= goal {
		streak++
		streakDays[day] = true
		day = day.AddDate(0, 0, -1)
	}
	return streak, todayMet, streakDays
}

// GetSummary builds the home-screen payload for the given date (default today).
// offsetMin is the client's getTimezoneOffset(), used so the streak is anchored
// to the user's local 3am-to-3am day.
func (s SummaryService) GetSummary(userID uint, dateStr string, offsetMin int) (*dto.SummaryResponse, *errors.APIError) {
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

	// Streak: consecutive local 3am-anchored days whose calories met the goal.
	entries, err := sm.AllEntryCalories(userID)
	if err != nil {
		return nil, errors.Internal("could not compute streak")
	}
	bucket := bucketCaloriesByStreakDay(entries, offsetMin)
	currentDay := localStreakDay(time.Now().UTC(), offsetMin)
	streak, todayMet, streakDays := computeCalorieStreak(bucket, calorieGoal, currentDay)

	// 7 streak-days ending on the current day. A day is a "hit" when its
	// calories met the goal; every day belonging to the active streak run also
	// carries the fire flag so the whole streak shows flames.
	week := make([]dto.DayState, 7)
	for i := range 7 {
		d := currentDay.AddDate(0, 0, i-6)
		met := calorieGoal > 0 && bucket[d] >= calorieGoal
		var state string
		switch {
		case d.Equal(currentDay):
			state = "today"
		case met:
			state = "hit"
		default:
			state = "miss"
		}
		week[i] = dto.DayState{Date: d.Format("2006-01-02"), Letter: dayLetters[d.Weekday()], State: state, Fire: streakDays[d]}
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
		Streak:            streak,
		StreakActiveToday: todayMet,
		Week:              week,
	}, nil
}

// GetStats builds the stats-screen payload for week / month / year ranges.
//
// week  → 7 daily data points
// month → 30 daily data points
// year  → 12 monthly aggregated data points
func (s SummaryService) GetStats(userID uint, rangeStr string, offsetMin int) (*dto.StatsResponse, *errors.APIError) {
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

	// Streak: consecutive local 3am-anchored days whose calories met the goal.
	entries, err := sm.AllEntryCalories(userID)
	if err != nil {
		return nil, errors.Internal("could not compute streak")
	}
	bucket := bucketCaloriesByStreakDay(entries, offsetMin)
	currentDay := localStreakDay(time.Now().UTC(), offsetMin)
	streak, _, _ := computeCalorieStreak(bucket, calorieGoal, currentDay)

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
