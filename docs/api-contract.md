# API Contract — Maxy Sports (Calorie Tracker)

**Owner:** Backend Dev maintains this file as the single source of truth.
**Consumer:** Frontend Dev builds repositories against it and requests changes via the team lead / Backend Dev.

This is a **seed proposal** derived from the Claude Design ("Maxy Sports Design System", project `3353971b-f6c2-4633-bd45-7cf9db15ca97`). Backend Dev may refine field shapes but must keep them aligned with what the screens need (below) and update this doc whenever the API changes. Frontend Dev: re-read this file before wiring each repository.

## Conventions

- Base URL: `http://localhost:3000/api/v1`
- Success envelope: `{ "data": <payload> }`  (existing `dto.BaseResponse`)
- Error envelope: `{ "message": <string> }` with appropriate HTTP status
- Auth: `Authorization: Bearer <token>` header on protected routes (JWT, 30-day)
- Timestamps: RFC3339 / ISO-8601 strings

## What the design screens need (requirements traceability)

- **LoginScreen** → Telegram login. On localhost the real Telegram widget cannot fire (domain/HTTPS restriction), so a **dev-login** path is required for testing.
- **OnboardingScreen** (3 steps) → captures `current_weight`, `goal_weight`, `direction (lose|maintain|gain)`, `calorie_goal`, and **optional** `carbs/fat/protein` goals → persist as the user's profile/goals.
- **HomeScreen** → today's totals (calories + macros), goals, current/goal weight, streak, last-7-days dot states; opens AddSheet and WeightSheet.
- **AddSheet** → `Quick add` (calories + optional macros), `New dish` (named, saved to history), `Recent` (re-add a previously named dish).
- **WeightSheet** → log today's weight; shows distance to goal.
- **StatsScreen** → calories-per-day series (week/month/year), weight trend series, avg/day, goal, days-under-goal, streak.

---

## Endpoints

### Auth
| Method | Path | Auth | Status |
|---|---|---|---|
| POST | `/auth/telegram` | public | **exists** |
| POST | `/auth/dev-login` | public (gated by `DEV_AUTH=true`) | **NEW** |
| GET  | `/auth/me` | bearer | **exists — extend** |

- `POST /auth/telegram` — body: Telegram widget payload `{id, first_name, last_name, username, photo_url, auth_date, hash}` → `{data:{token}}`.
- `POST /auth/dev-login` — **dev-only** (return 404 when `DEV_AUTH` is not `true`). Body `{ "username"?: string }`. Upserts a deterministic dev user so the full flow is testable on localhost without a real Telegram domain → `{data:{token}}`.
- `GET /auth/me` — extend response to include `onboarded: bool` so the frontend can route new users to onboarding:
  `{data:{ id, telegram_id, first_name, last_name, username, photo_url, created_at, onboarded }}`.

### Profile / Goals  (NEW domain — mirror the food domain: schema → model → service → controller → routes → dto)
| Method | Path | Auth |
|---|---|---|
| GET | `/profile` | bearer |
| PUT | `/profile` | bearer |

- `GET /profile` → `{data:{ calorie_goal, carbs_goal, fat_goal, protein_goal, current_weight, goal_weight, direction, onboarded }}`. Return zero/default values with `onboarded:false` when not yet set.
- `PUT /profile` — body `{ calorie_goal, carbs_goal?, fat_goal?, protein_goal?, current_weight?, goal_weight?, direction? }`. Used by onboarding "Start tracking" and later edits. Sets `onboarded:true`. `direction ∈ lose|maintain|gain`. When `current_weight` is supplied, also seed an initial weight entry so the trend has a starting point.

### Weight  (NEW domain)
| Method | Path | Auth |
|---|---|---|
| POST | `/weights` | bearer |
| GET  | `/weights` | bearer |

- `POST /weights` — body `{ weight: number, recorded_at?: string }` (defaults to now) → `{data:{ id, weight, recorded_at }}`.
- `GET /weights` → `{data:[{ id, weight, recorded_at }]}` (chronological). Current weight = latest entry.

### Food  (exists + extend)
| Method | Path | Auth | Status |
|---|---|---|---|
| POST | `/foods` | bearer | **exists** |
| GET  | `/foods?date=YYYY-MM-DD` | bearer | **exists — add optional `date` filter** |
| DELETE | `/foods/:id` | bearer | **exists** |
| GET | `/foods/recent` | bearer | **NEW** |

- `GET /foods/recent` → distinct recent **named** dishes (exclude unnamed "Quick add"), most-recent-first, limited (~20): `{data:[{ name, calories, protein, carbs, fat }]}`. Powers AddSheet "Recent".

### Summary & Stats  (NEW domain)
| Method | Path | Auth |
|---|---|---|
| GET | `/summary?date=YYYY-MM-DD` | bearer |
| GET | `/stats?range=week\|month\|year` | bearer |

- `GET /summary` (date optional, default today) — one call powering Home:
  ```json
  { "data": {
      "date": "2026-06-17",
      "totals": { "calories": 1240, "carbs": 124, "fat": 38, "protein": 90 },
      "goals":  { "calories": 2400, "carbs": 270, "fat": 80, "protein": 150 },
      "weight": { "current": 84.2, "goal": 78 },
      "streak": 12,
      "week": [ { "date": "2026-06-15", "letter": "S", "state": "hit" } ]
  } }
  ```
  `week` = 7 entries; `state ∈ hit|miss|today|future`.
- `GET /stats` — powers Stats:
  ```json
  { "data": {
      "range": "week",
      "calories_per_day": [ { "label": "M", "date": "2026-06-09", "value": 2280, "muted": false } ],
      "weight_trend":     [ { "label": "Wk1", "value": 84.2, "recorded_at": "..." } ],
      "avg_per_day": 2400,
      "goal": 2400,
      "days_under_goal": 5,
      "streak": 12
  } }
  ```

---

## Change log
- _(seed)_ Initial proposal from design + existing backend. Backend Dev: append entries here as you implement/adjust.
- _(2026-06-17)_ **Implemented all endpoints.** Dockerfile upgraded to golang:1.26-alpine. DEV_AUTH env var added. New schemas: `user_profiles`, `weight_entries`. New domains: profile, weight, summary/stats. Food extended: optional `?date=YYYY-MM-DD` filter, `GET /foods/recent`. Auth extended: `POST /auth/dev-login` (gated by DEV_AUTH=true), `GET /auth/me` now returns `onboarded: bool`. Stats year range returns 12 monthly-aggregated data points. Server running on localhost:3000.
