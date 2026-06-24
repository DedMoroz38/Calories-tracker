# Test Plan — Maxy Sports Calorie Tracker

**Owner:** QA  
**Grounded in:** `docs/api-contract.md`, `CLAUDE.md`, existing backend code  
**Phase 1 status:** Written while devs build.  
**Phase 2:** Execute once BackendDev + FrontendDev both report ready.

---

## Scope

This plan covers backend (Go/Fiber on :3000), frontend (Next.js on :3001), and their integration. It is structured so each check can produce an unambiguous PASS or FAIL for `tests/report.md`.

---

## 1. Backend Checks

### 1.1 Build

| ID | Check | Command | Pass Criteria |
|----|-------|---------|---------------|
| B-BUILD-1 | Binary compiles | `go build ./cmd/api` | Exit 0, no errors |
| B-BUILD-2 | All tests pass | `go test ./...` | Exit 0, no failures |

### 1.2 Auth — POST /api/v1/auth/telegram

| ID | Check | Method | Pass Criteria |
|----|-------|--------|---------------|
| B-AUTH-1 | Valid Telegram payload returns token | POST with correctly signed HMAC body | 200, `{"data":{"token":"<jwt>"}}` |
| B-AUTH-2 | Missing `id` field rejected | POST `{}` | 400, `{"message":"..."}` |
| B-AUTH-3 | Missing `hash` rejected | POST without hash | 400 |
| B-AUTH-4 | Missing `auth_date` rejected | POST without auth_date | 400 |
| B-AUTH-5 | Invalid hash rejected | POST with bad hash value | 401 |
| B-AUTH-6 | Expired auth_date rejected | POST with auth_date > 24 h ago | 401 |
| B-AUTH-7 | Response envelope correct | Valid request | Exactly `{"data":{"token":"..."}}`, no extra fields |

### 1.3 Auth — POST /api/v1/auth/dev-login (NEW)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-DEV-1 | Returns 404 when DEV_AUTH is not set/false | 404 |
| B-DEV-2 | Returns 404 when DEV_AUTH=false | 404 |
| B-DEV-3 | Returns 200 + JWT when DEV_AUTH=true | 200, `{"data":{"token":"..."}}` |
| B-DEV-4 | Optional `username` in body accepted | 200 with username present |
| B-DEV-5 | Omitted `username` accepted (deterministic dev user) | 200 |
| B-DEV-6 | Same token-issuing user upserted on repeat calls | Second call returns valid JWT for same user |

### 1.4 Auth — GET /api/v1/auth/me

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-ME-1 | No token returns 401 | 401 |
| B-ME-2 | Invalid/expired token returns 401 | 401 |
| B-ME-3 | Valid token returns user | 200, `{"data":{...}}` |
| B-ME-4 | Response includes `onboarded` field | `data.onboarded` is bool (true or false) |
| B-ME-5 | Response shape is exact | Fields: `id, telegram_id, first_name, last_name, username, photo_url, created_at, onboarded` — no extra fields leaked |

### 1.5 Profile — GET /api/v1/profile (NEW)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-PROF-1 | No token returns 401 | 401 |
| B-PROF-2 | New user returns zeros + onboarded:false | 200, all numeric goals are 0/null, `onboarded:false` |
| B-PROF-3 | Response shape correct | Fields: `calorie_goal, carbs_goal, fat_goal, protein_goal, current_weight, goal_weight, direction, onboarded` |
| B-PROF-4 | After PUT, GET returns updated values | Subsequent GET reflects PUT body |

### 1.6 Profile — PUT /api/v1/profile (NEW)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-PROF-5 | No token returns 401 | 401 |
| B-PROF-6 | Valid body sets profile | 200, profile saved |
| B-PROF-7 | Sets onboarded:true | After PUT, GET /profile returns `onboarded:true` |
| B-PROF-8 | direction must be lose/maintain/gain | Invalid direction → 400 |
| B-PROF-9 | carbs/fat/protein are optional | PUT with only calorie_goal → 200 |
| B-PROF-10 | Seeding current_weight also creates weight entry | After PUT with current_weight, GET /weights returns that entry |
| B-PROF-11 | Ownership: user A cannot GET/PUT user B's profile | Requests with user A's JWT only affect user A's data |

### 1.7 Weights — POST /api/v1/weights (NEW)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-WT-1 | No token returns 401 | 401 |
| B-WT-2 | Valid body logs weight | 201 (or 200), `{"data":{"id":...,"weight":...,"recorded_at":"..."}}` |
| B-WT-3 | Missing `weight` field rejected | 400 |
| B-WT-4 | Optional `recorded_at` defaults to now | Omitting recorded_at → entry created with ~now |
| B-WT-5 | Response shape exact | Fields: `id, weight, recorded_at` only |

### 1.8 Weights — GET /api/v1/weights (NEW)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-WT-6 | No token returns 401 | 401 |
| B-WT-7 | Returns array of weight entries | 200, `{"data":[...]}` |
| B-WT-8 | Entries are chronological | Oldest first |
| B-WT-9 | Ownership: only own entries returned | User A cannot see user B's weights |

### 1.9 Food — POST /api/v1/foods

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-FOOD-1 | No token returns 401 | 401 |
| B-FOOD-2 | Valid body creates entry | 201, `{"data":{...}}` |
| B-FOOD-3 | Missing `name` rejected | 400 |
| B-FOOD-4 | Negative `calories` rejected | 400 |
| B-FOOD-5 | Negative macros rejected | 400 |
| B-FOOD-6 | Missing `consumed_at` defaults to now | Entry has plausible timestamp |

### 1.10 Food — GET /api/v1/foods

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-FOOD-7 | No token returns 401 | 401 |
| B-FOOD-8 | Returns list of entries | 200, `{"data":[...]}` |
| B-FOOD-9 | `date=YYYY-MM-DD` filter works | Only entries from that date returned |
| B-FOOD-10 | Wrong date format returns 400 | 400 |
| B-FOOD-11 | Ownership: only own entries returned | User A's JWT cannot return user B's food |

### 1.11 Food — DELETE /api/v1/foods/:id

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-FOOD-12 | No token returns 401 | 401 |
| B-FOOD-13 | Valid id deletes entry | 200, deletion confirmed |
| B-FOOD-14 | Non-existent id returns 404 | 404 |
| B-FOOD-15 | Ownership: cannot delete another user's entry | DELETE with user A's JWT against user B's food id → 404 |

### 1.12 Food — GET /api/v1/foods/recent (NEW)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-FOOD-16 | No token returns 401 | 401 |
| B-FOOD-17 | Returns distinct named dishes | 200, `{"data":[...]}` |
| B-FOOD-18 | Quick-add entries (unnamed) excluded | Entries with blank name not in results |
| B-FOOD-19 | Most-recent-first ordering | Verify order by created_at DESC |
| B-FOOD-20 | Max ~20 results | Response length <= 20 |
| B-FOOD-21 | Response shape: name, calories, protein, carbs, fat only | No id, user_id, consumed_at leaked |

### 1.13 Summary — GET /api/v1/summary (NEW)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-SUM-1 | No token returns 401 | 401 |
| B-SUM-2 | No date param defaults to today | `data.date` equals today in YYYY-MM-DD |
| B-SUM-3 | `date=YYYY-MM-DD` returns that date | `data.date` matches param |
| B-SUM-4 | Totals reflect actual food entries | Add food then verify totals |
| B-SUM-5 | Goals match profile | Goals match what was set in PUT /profile |
| B-SUM-6 | Weight present when weight logged | `data.weight.current` reflects latest weight entry |
| B-SUM-7 | Streak calculation is correct | Log food for N days consecutively; streak = N |
| B-SUM-8 | `week` array has 7 entries | `data.week.length === 7` |
| B-SUM-9 | `week[].state` values are valid | Each state ∈ `hit|miss|today|future` |
| B-SUM-10 | Ownership: scoped to requesting user | Data not cross-contaminated |
| B-SUM-11 | Response envelope shape | Exactly the fields in api-contract.md; no extras |

### 1.14 Stats — GET /api/v1/stats (NEW)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| B-STATS-1 | No token returns 401 | 401 |
| B-STATS-2 | `range=week` returns 7-day series | 200, calories_per_day has 7 entries |
| B-STATS-3 | `range=month` returns month series | 200 |
| B-STATS-4 | `range=year` returns year series | 200 |
| B-STATS-5 | Missing/invalid range returns 400 | 400 |
| B-STATS-6 | `avg_per_day` is arithmetic average | Verify against known food data |
| B-STATS-7 | `days_under_goal` count is correct | Count of days where calories < calorie_goal |
| B-STATS-8 | `streak` matches summary streak | Same value as /summary |
| B-STATS-9 | Response envelope shape | Exactly the fields in api-contract.md |

---

## 2. Frontend Checks

### 2.1 Build & Type Safety

| ID | Check | Command | Pass Criteria |
|----|-------|---------|---------------|
| F-BUILD-1 | `next build` succeeds | `cd frontend && npm run build` | Exit 0, no errors |
| F-BUILD-2 | TypeScript clean | `npx tsc --noEmit` | Exit 0, 0 errors |
| F-BUILD-3 | No `any` leaks in repo/DTO layers | grep for `: any` / `as any` in src/ | Zero hits in repository and DTO files |

### 2.2 Project Structure (FSD)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| F-FSD-1 | Feature Sliced Design structure respected | Top-level: `app/`, `pages/` (or `views/`), `features/`, `entities/`, `shared/` |
| F-FSD-2 | Single shared httpClient | One httpClient module in `shared/`; no other fetch/axios usages |
| F-FSD-3 | Each screen wired to its own repository | Repositories in `entities/` or `features/`; each calls the correct endpoint |
| F-FSD-4 | No direct fetch() calls in components | All API calls go through repository layer |

### 2.3 DTOs / API Alignment

| ID | Check | Pass Criteria |
|----|-------|---------------|
| F-DTO-1 | Auth DTO matches contract | Sends exactly `{id, first_name, last_name, username, photo_url, auth_date, hash}` |
| F-DTO-2 | Dev-login DTO minimal | Body is `{ username?: string }` only |
| F-DTO-3 | Profile PUT DTO minimal | Only fields the API requires; no extras |
| F-DTO-4 | Food create DTO minimal | `{name, calories, protein?, carbs?, fat?, consumed_at?}` |
| F-DTO-5 | Weight create DTO minimal | `{weight, recorded_at?}` |
| F-DTO-6 | No extraneous fields sent | Inspect network payloads: no junk fields in request bodies |
| F-DTO-7 | TypeScript types match response shapes | Response types in repos carry only fields the API returns |

### 2.4 Token Handling

| ID | Check | Pass Criteria |
|----|-------|---------------|
| F-TOKEN-1 | JWT stored on login | After dev-login, token in localStorage (or cookie — consistent with backend) |
| F-TOKEN-2 | Token attached to every protected request | Authorization: Bearer header present on /foods, /profile, etc. |
| F-TOKEN-3 | Missing token redirects to login | Navigating protected page without token shows login screen |
| F-TOKEN-4 | 401 response redirects to login | Any 401 from API causes logout/redirect |

### 2.5 Onboarding Routing

| ID | Check | Pass Criteria |
|----|-------|---------------|
| F-ONB-1 | New user (onboarded:false) routed to onboarding | After dev-login with fresh user, sees OnboardingScreen |
| F-ONB-2 | Returning user (onboarded:true) routed to home | After completing onboarding, direct to HomeScreen |
| F-ONB-3 | Onboarding 3-step flow works | All 3 steps navigate correctly |
| F-ONB-4 | "Start tracking" submits PUT /profile | Network call to /profile on final step |
| F-ONB-5 | After onboarding, back-nav does not re-show onboarding | Router guards work |

### 2.6 HomeScreen

| ID | Check | Pass Criteria |
|----|-------|---------------|
| F-HOME-1 | Loads GET /summary | Network call to /summary on mount |
| F-HOME-2 | Totals displayed correctly | Calories + macro totals match summary response |
| F-HOME-3 | Goals displayed correctly | Goal values match profile |
| F-HOME-4 | Streak shown | Streak number visible |
| F-HOME-5 | 7-day week dots rendered | All 7 dots shown with correct hit/miss/today/future state |
| F-HOME-6 | Loading state shown while fetching | Spinner or skeleton during request |
| F-HOME-7 | Error state shown on failure | Visible error if /summary call fails |

### 2.7 AddSheet (Quick add / New dish / Recent)

| ID | Check | Pass Criteria |
|----|-------|---------------|
| F-ADD-1 | Quick-add submits POST /foods | Network call with calories (+ optional macros) |
| F-ADD-2 | New dish submits POST /foods with name | name field populated |
| F-ADD-3 | Recent tab loads GET /foods/recent | Network call on tab open |
| F-ADD-4 | Re-adding a recent dish submits POST /foods | POST with same name/macros |
| F-ADD-5 | After add, summary refreshes | Home totals update without full reload |

### 2.8 WeightSheet

| ID | Check | Pass Criteria |
|----|-------|---------------|
| F-WT-1 | Submits POST /weights | Network call with weight value |
| F-WT-2 | Distance-to-goal displayed | Shows `current - goal` with correct sign |
| F-WT-3 | Error state on failure | Network error is surfaced |

### 2.9 StatsScreen

| ID | Check | Pass Criteria |
|----|-------|---------------|
| F-STATS-1 | Loads GET /stats?range=week on mount | Network call |
| F-STATS-2 | Range selector changes query param | Switching to month/year triggers new /stats call |
| F-STATS-3 | Calories chart renders | Chart visible with correct labels |
| F-STATS-4 | Weight trend chart renders | Chart visible |
| F-STATS-5 | avg_per_day, days_under_goal, streak displayed | Numbers shown in UI |

---

## 3. Integration Checks

| ID | Check | Method | Pass Criteria |
|----|-------|--------|---------------|
| I-1 | Backend starts and db connects | `docker compose up db -d && go run ./cmd/api` | Server listens on :3000, no panic |
| I-2 | Frontend starts | `npm run dev` in frontend/ | Next.js dev server on :3001 |
| I-3 | Dev-login end-to-end | Browser → /auth/dev-login via UI | User logged in, token in storage |
| I-4 | Onboarding end-to-end | Complete 3-step onboarding | /profile updated, onboarded:true, routed to home |
| I-5 | Add food end-to-end | Add food via AddSheet | Appears in today's food log, totals update |
| I-6 | Log weight end-to-end | Log weight via WeightSheet | New weight entry in /weights |
| I-7 | Stats reflect logged data | Visit StatsScreen after adding data | Charts show the data just logged |
| I-8 | CORS not blocking frontend | All API calls from :3001 to :3000 succeed | No CORS errors in browser console |
| I-9 | Delete food end-to-end | Delete a food entry | Entry removed, totals update |

---

## 4. Known Gaps to Track (before Phase 2)

The following are expected to be missing from the current backend; Phase 2 will verify they are implemented:

- `POST /auth/dev-login` with DEV_AUTH gating
- `GET /auth/me` `onboarded` field
- `GET|PUT /profile` domain
- `POST|GET /weights` domain
- `GET /foods?date=` filter
- `GET /foods/recent`
- `GET /summary` with streak + week logic
- `GET /stats` with range param
- CORS middleware on backend

---

## 5. Phase 2 Execution Steps

When both devs report ready:

1. Run `go test ./...` — verify B-BUILD-2.
2. Run `go build ./cmd/api` — verify B-BUILD-1.
3. Start the stack: `docker compose up db -d && go run ./cmd/api`.
4. Run `curl` / `httpie` tests against each B-* check systematically.
5. Run `cd frontend && npm run build` — verify F-BUILD-1.
6. Run `npx tsc --noEmit` — verify F-BUILD-2.
7. Start `npm run dev`.
8. Open Playwright browser to http://localhost:3001 and drive through I-3 → I-9.
9. Inspect network calls and console for F-DTO-*, F-TOKEN-*, F-FSD-* checks.
10. Write `tests/report.md` with PASS/FAIL for every check above.
11. Send defect list to BackendDev/FrontendDev.
12. Repeat until all checks pass, then send QA APPROVED to main.
