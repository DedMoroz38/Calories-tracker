# Investigation: "Today's log" shows entries from other days

Status: **investigated, not yet fixed** (2026-06-25)

## Symptom

The home screen's **Today's log** card lists food entries from previous days,
not just today. Meanwhile the calorie ring / macro bars at the top of the same
card only reflect *today's* totals — so the two halves of the screen disagree.

## Root cause

The home view requests the food log **without a date filter**, so the backend
returns *all* of the user's entries across all days.

### Data path

1. `frontend/views/home/index.tsx` (`loadData`) calls:
   ```ts
   getFoods()   // ← no argument
   ```
2. `frontend/entities/food/api/foodRepository.ts`:
   ```ts
   export async function getFoods(date?: string) {
     const qs = date ? `?date=${date}` : "";   // no date → no query string
     return http.get<FoodEntry[]>(`/foods${qs}`);
   }
   ```
3. Backend `FoodController.List` (`backend/internal/MVC/controllers/food_controller.go`)
   only filters when `?date=` is present; otherwise `day` stays the zero time.
4. `FoodModel.ListByUser` (`backend/internal/MVC/models/food_model.go`):
   ```go
   if !day.IsZero() {          // zero day → this block is skipped
       q = q.Where("consumed_at >= ? AND consumed_at < ?", start, end)
   }
   ```
   With a zero `day`, **no date predicate is applied → every entry is returned.**

### Why the ring is correct but the log isn't

The summary endpoint (`GET /summary`) defaults to *today* when no date is given
(`summary_service.go`: `target = time.Now().UTC()`, then `TotalsForDay`). So the
ring/macros are scoped to today while the log is not — producing the mismatch.

## Fix options

- **Preferred (frontend):** pass today's date to `getFoods` from the home view,
  e.g. `getFoods(todayYYYYMMDD)`. Smallest change, keeps the "all entries"
  capability available for other callers (e.g. a history view).
- **Alternative (backend):** make `List` default to today when no `?date=` is
  supplied. Changes the endpoint's default semantics; would break any caller
  that relies on getting the full history.

## Secondary concern: UTC vs local calendar day

The day boundary is computed in **UTC** everywhere on the backend
(`ListByUser`, `TotalsForDay`, streak), but the header label uses the browser's
**local** date (`new Date().toLocaleDateString(...)`). For users far from UTC,
"today" on the device can differ from "today" in the DB query — entries logged
late at night (local) may land on the previous/next UTC day, causing off-by-one
inclusion/exclusion even after the primary fix. Decide whether the date filter
should be local-day based (pass the client's local `YYYY-MM-DD`, which the
preferred frontend fix naturally does) and whether the summary/streak should
likewise be made timezone-aware.

## Relevant files

- `frontend/views/home/index.tsx` — `loadData`, calls `getFoods()`
- `frontend/entities/food/api/foodRepository.ts` — `getFoods(date?)`
- `backend/internal/MVC/controllers/food_controller.go` — `List`
- `backend/internal/MVC/models/food_model.go` — `ListByUser`
- `backend/internal/MVC/services/summary_service.go` — today-default for summary
