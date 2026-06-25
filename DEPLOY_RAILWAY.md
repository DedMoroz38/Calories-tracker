# Deploying Maxy Sports on Railway

This guide deploys the **whole app** — Go backend, Next.js frontend, PostgreSQL,
and S3-compatible object storage — to [Railway](https://railway.app), then puts
it behind your custom domain.

Railway has no native S3 service, so we use the **MinIO** template, which speaks
the S3 protocol. The backend already supports any S3-compatible endpoint via the
`S3_ENDPOINT` variable — no code changes are needed to point it at MinIO.

---

## Architecture

One Railway **project** holds four **services**:

| Service    | What it is                     | Source                         |
|------------|--------------------------------|--------------------------------|
| `postgres` | PostgreSQL database            | Railway database plugin        |
| `minio`    | S3-compatible object storage   | Railway "MinIO" template       |
| `backend`  | Go + Fiber API                 | this repo, root dir `backend/` |
| `frontend` | Next.js app                    | this repo, root dir `frontend/`|

Photos and avatars are stored **privately** in MinIO and served to the browser
through short-lived **presigned URLs** the backend generates. The MinIO bucket
never needs to be public — but its API **does** need a public Railway domain so
those presigned URLs are reachable from users' browsers (covered below).

---

## Prerequisites

- A Railway account and the repo pushed to GitHub (Railway deploys from GitHub).
- Your custom domain's DNS managed somewhere you can add `CNAME` records.
- Your Telegram bot token (already in `backend/.env`).

---

## Step 1 — Create the project and PostgreSQL

1. Railway dashboard → **New Project** → **Deploy PostgreSQL**.
2. This creates the `postgres` service. Open it → **Variables** and note that it
   exposes `DATABASE_URL` (you'll reference it, not copy it).

---

## Step 2 — Add MinIO (object storage)

1. In the same project → **New** → **Template** → search **MinIO** → deploy.
2. Open the `minio` service → **Variables**. Note (or set) these:
   - `MINIO_ROOT_USER` — the admin username (your S3 access key).
   - `MINIO_ROOT_PASSWORD` — the admin password (your S3 secret key).
3. **Expose the MinIO API publicly** (needed so browsers can load presigned URLs):
   - `minio` service → **Settings** → **Networking** → **Generate Domain**.
   - MinIO serves its **S3 API on port 9000** and its web console on **9001**.
     Generate a public domain bound to **port 9000**. This URL is your
     `S3_ENDPOINT`, e.g. `https://minio-production-abcd.up.railway.app`.
   - (Optional) Generate a second domain on port 9001 to reach the web console.
4. **Pick a bucket name** (e.g. `maxy-photos`). You do **not** need to create it
   by hand — the backend creates the bucket automatically on first boot using
   your credentials. Just use this name as `AWS_S3_BUCKET` in Step 3.

> You can use the root user/password directly as the access/secret key, or
> create a dedicated **Access Key** in the console (Identity → Access Keys) and
> use that instead — recommended for production.

---

## Step 3 — Deploy the backend

1. Project → **New** → **GitHub Repo** → select this repo.
2. Open the new service → **Settings**:
   - **Root Directory**: `backend`
   - Railway auto-detects the existing `backend/Dockerfile` and builds it.
   - Leave the start command empty (the Dockerfile's `CMD ["./server"]` runs it).
3. **Variables** (Settings → Variables) — set:

   | Variable                | Value                                                        |
   |-------------------------|--------------------------------------------------------------|
   | `DATABASE_URL`          | `${{Postgres.DATABASE_URL}}` (reference the pg service)       |
   | `JWT_SECRET`            | a long random string                                         |
   | `TELEGRAM_BOT_TOKEN`    | your bot token                                               |
   | `S3_ENDPOINT`           | the **port-9000** MinIO public URL from Step 2               |
   | `AWS_S3_BUCKET`         | your bucket name, e.g. `maxy-photos`                         |
   | `AWS_ACCESS_KEY_ID`     | `${{minio.MINIO_ROOT_USER}}` (or your access key)            |
   | `AWS_SECRET_ACCESS_KEY` | `${{minio.MINIO_ROOT_PASSWORD}}` (or your secret key)        |
   | `AWS_REGION`            | leave unset (defaults to `us-east-1`) — MinIO ignores it     |

   > Do **not** set `PORT` — Railway injects it and the Go server reads it
   > automatically (it listens on `:$PORT`).

4. **Settings → Networking → Generate Domain** for the backend. You'll get
   something like `https://backend-production-1234.up.railway.app`. Note it —
   the frontend needs it next.
5. Deploy. Check **Deploy Logs**: you should **not** see the
   `WARNING: S3 not configured` line. If you do, a storage var is missing.

---

## Step 4 — Deploy the frontend

1. Project → **New** → **GitHub Repo** → same repo again.
2. Open the service → **Settings**:
   - **Root Directory**: `frontend`
   - Railway auto-detects Next.js (Nixpacks): it runs `npm install`,
     `npm run build`, then `npm start` (which is `next start` and binds `$PORT`).
3. **Variables** — set the API base URL the browser will call:

   | Variable               | Value                                          |
   |------------------------|------------------------------------------------|
   | `NEXT_PUBLIC_API_BASE` | your **backend** public URL from Step 3         |

   > ⚠️ **This is a build-time variable.** `NEXT_PUBLIC_*` values are inlined
   > into the client bundle when `next build` runs — they are **not** read at
   > runtime. So whenever you change `NEXT_PUBLIC_API_BASE` (e.g. after adding a
   > custom domain in Step 5), you must **redeploy the frontend** for it to take
   > effect.

4. **Settings → Networking → Generate Domain** for the frontend (or skip
   straight to the custom domain in Step 5).
5. Deploy.

At this point the app is live on the generated Railway domains.

---

## Step 5 — Point your custom domain at the app

You bought a domain; here's how to attach it. Decide your layout, e.g.:

- `app.yourdomain.com` → **frontend**
- `api.yourdomain.com` → **backend**

For **each** service:

1. Service → **Settings → Networking → Custom Domain** → enter the subdomain
   (e.g. `app.yourdomain.com`).
2. Railway shows a **CNAME target** (e.g. `xyz.up.railway.app`).
3. In your DNS provider, add a `CNAME` record:
   `app` → the target Railway gave you. Repeat for `api` → backend's target.
4. Wait for DNS to propagate; Railway provisions TLS automatically (the domain
   shows a green check when ready).

> Apex domains (`yourdomain.com` with no subdomain) can't be a CNAME on most DNS
> providers — use a subdomain like `app.`, or use your provider's ALIAS/flattened
> CNAME feature if it has one.

---

## Step 6 — Re-point the frontend at the backend's custom domain

Now that the backend has `https://api.yourdomain.com`, update the frontend so the
browser calls the nice domain instead of the `*.up.railway.app` one:

1. `frontend` service → **Variables** → set
   `NEXT_PUBLIC_API_BASE = https://api.yourdomain.com`
2. **Redeploy the frontend** (Deployments → ⋮ → Redeploy). This is required —
   see the build-time note in Step 4. The value is baked into the JS bundle.

**Where the base URL lives in code:** `frontend/shared/config/index.ts`:

```ts
export const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:3000";
export const API_V1 = `${API_BASE}/api/v1`;
```

So you never edit code to change environments — only the
`NEXT_PUBLIC_API_BASE` variable (then redeploy). Locally it falls back to
`http://localhost:3000`.

---

## Step 7 — Update Telegram login domain

The Telegram Login Widget only renders on the domain registered with BotFather.

1. In Telegram, message **@BotFather** → `/setdomain` → pick your bot
   (`MaxyCalorieTrackerBot`).
2. Send your **frontend** domain **without scheme or path**:
   `app.yourdomain.com`
3. Reload the app at `https://app.yourdomain.com` and hard-refresh — the login
   button should render.

---

## Environment variable reference

### backend
| Variable                | Required | Notes                                            |
|-------------------------|----------|--------------------------------------------------|
| `DATABASE_URL`          | yes      | reference `${{Postgres.DATABASE_URL}}`           |
| `JWT_SECRET`            | yes      | long random string                               |
| `TELEGRAM_BOT_TOKEN`    | yes      | from BotFather                                   |
| `S3_ENDPOINT`           | yes*     | MinIO port-9000 public URL (*omit for real AWS)  |
| `AWS_S3_BUCKET`         | yes      | bucket name                                      |
| `AWS_ACCESS_KEY_ID`     | yes      | MinIO user / access key                          |
| `AWS_SECRET_ACCESS_KEY` | yes      | MinIO password / secret key                      |
| `AWS_REGION`            | no       | defaults to `us-east-1`                          |
| `PORT`                  | no       | injected by Railway                              |

### frontend
| Variable               | Required | Notes                                             |
|------------------------|----------|---------------------------------------------------|
| `NEXT_PUBLIC_API_BASE` | yes      | backend public URL; **build-time**, redeploy on change |
| `PORT`                 | no       | injected by Railway                               |

---

## Switching back to real AWS S3 later

The code is provider-agnostic. To use genuine AWS S3 instead of MinIO: leave
`S3_ENDPOINT` **unset**, and set `AWS_REGION`, `AWS_S3_BUCKET`,
`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` to your AWS values. Path-style
addressing turns off automatically when `S3_ENDPOINT` is empty.

---

## Troubleshooting

- **`WARNING: S3 not configured` in backend logs** — one of `S3_ENDPOINT` (for
  MinIO), `AWS_S3_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` is
  missing. (`S3_ENDPOINT` is only required for MinIO; on AWS it's expected to be
  empty.)
- **Photos upload but won't display / 403 on image URLs** — the browser can't
  reach `S3_ENDPOINT`. Make sure it's the **public, port-9000** MinIO domain
  (not the internal `*.railway.internal` host, and not the 9001 console).
- **Feed/profile return `"photo storage is not configured"`** — same as the
  warning above; the backend booted without storage vars.
- **Login button missing** — `/setdomain` in BotFather must exactly match the
  frontend domain, served over HTTPS.
- **Frontend still calls the old API URL** — `NEXT_PUBLIC_API_BASE` is baked at
  build time; redeploy the frontend after changing it.
- **`ensure bucket … failed` on backend startup** — the backend couldn't create
  the bucket. Usually means the credentials lack permission, or `S3_ENDPOINT` is
  wrong/unreachable. Root MinIO credentials always have permission; a scoped
  access key needs `s3:CreateBucket` (or pre-create the bucket once in the
  console and the backend will just use it).
