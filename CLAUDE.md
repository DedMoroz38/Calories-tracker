# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Telegram-authenticated calorie tracker. The backend is a Go REST API; the frontend is a single static HTML file served separately.

## Commands

All backend commands run from `backend/`.

```bash
# Start Postgres (required before running the server)
docker compose up db -d

# Run the server locally
go run ./cmd/api

# Build
go build -o server ./cmd/api

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/auth/...

# Run the full stack via Docker
docker compose up --build
```

The frontend (`frontend/index.html`) is a plain HTML file — open it directly in a browser or serve it with any static server:

```bash
python3 -m http.server 8080 --directory frontend
```

## Architecture

**Backend** (`backend/`) — Go + Fiber v2, GORM, PostgreSQL.

- `cmd/api/main.go` — entry point: wires config → DB → Fiber routes → handlers.
- `internal/config/` — loads env vars (`.env` via godotenv). Required vars: `PORT`, `DATABASE_URL`, `TELEGRAM_BOT_TOKEN`, `JWT_SECRET`.
- `internal/auth/telegram.go` — verifies the Telegram Login Widget payload (HMAC-SHA256 over sorted `key=value` pairs, 24 h expiry window).
- `internal/auth/jwt.go` — issues/parses HS256 JWTs (30-day expiry).
- `internal/handlers/auth.go` — `POST /auth/telegram`: verifies Telegram data, upserts `User` via GORM `FirstOrCreate`, returns JWT.
- `internal/middleware/auth.go` — Bearer-token middleware for protected routes; sets `c.Locals("userID")`.
- `internal/models/user.go` — GORM `User` model; `AutoMigrate` runs on startup.
- `internal/dto/food.go` — request DTO for food entries (handler not yet wired).

**Frontend** (`frontend/index.html`) — single file, no build step.

- Embeds the official Telegram Login Widget pointing at `MaxyCalorieTrackerBot`.
- On auth callback `onTelegramAuth(user)`, POSTs the raw widget payload to `${API_BASE}/auth/telegram` and stores the returned JWT in `localStorage`.
- `API_BASE` is hardcoded to `http://127.0.0.1:3000` — change it when deploying.

## Key Constraints

- **Telegram widget domain restriction**: The Login Widget only fires on pages served over HTTPS from a domain registered with BotFather (`/setdomain`). `file://` and `http://localhost` will not work in production; use ngrok or a real domain during development.
- **CORS**: The backend has no CORS middleware. Add `github.com/gofiber/contrib/cors` (or `fiber/v2`'s built-in) before the frontend and backend are on different origins.
- **No CORS / no HTTPS locally**: The widget will silently fail to render if the serving domain doesn't match the bot's allowed domain.
