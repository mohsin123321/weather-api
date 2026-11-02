# Weather API — Implementation Plan

## Problem

Weather data is useful but providers are unreliable — keys expire, APIs go down, rate limits hit. A single-provider proxy with no resilience, no observability, and a static HTML frontend is not production-ready.

**Goal:** evolve this into a full-stack application that demonstrates real backend design — multiple providers with fallback, Redis caching, observability, rate limiting — paired with a modern React frontend and a live deployment.

---

## Architecture

```
Browser (React + Vite)
        │  dev:  :5173 → proxy → :8080
        │  prod: Go serves web/dist/
        ▼
  cmd/api  (Go HTTP server)
        │
   ┌────┴────┐
   │  cache  │  in-memory now → Redis (Step 5)
   └────┬────┘
        │
  provider router  ←── circuit breaker (Step 4)
   ┌────┴──────┐
   │           │
  OWM      Open-Meteo
```

---

## Directory Layout

```
/
├── cmd/api/              # entry point
├── internal/
│   ├── adapters/         # weather provider adapters (OWM, Open-Meteo)
│   ├── apperrors/        # typed errors with HTTP status codes
│   ├── config/           # env + .env config loading
│   ├── domain/           # port interfaces (adapter, cache)
│   ├── dto/              # response types mirroring provider APIs
│   ├── persistence/      # cache implementations (in-memory → Redis)
│   └── server/           # HTTP handlers, middleware, routes
├── web/                  # React + TypeScript (Vite)
│   ├── src/
│   │   ├── api/          # typed fetch wrappers
│   │   ├── components/   # WeatherCard, ForecastChart, CityComparison
│   │   └── types/        # TypeScript types mirroring Go DTOs
│   ├── Dockerfile        # Vite dev server container
│   └── dist/             # production build, served by Go
├── tests/                # mock generators and test fixtures
├── Dockerfile            # production: Node build → Go build → alpine
├── Dockerfile.api        # dev: Go API only (used by docker-compose)
├── docker-compose.yml    # api + web services
└── .github/workflows/    # CI pipeline (Step 3)
```

---

## Progress

| Step | What | Status |
|---|---|---|
| 1 | Foundation: Go API, OWM adapter, in-memory cache, chi router, CORS, Swagger, unit tests, graceful shutdown | done |
| 2 | React frontend: Vite + React + TypeScript, WeatherCard component, Vite proxy, Docker split (api + web services) | done |
| 3 | Docker + CI: production Dockerfile, GitHub Actions (test → lint → build) | in progress |
| 4 | Multiple providers + circuit breaker: Open-Meteo fallback, three-state circuit breaker | — |
| 5 | Redis cache: replace in-memory, stale-while-revalidate, docker-compose Redis service | — |
| 6 | Observability: Prometheus metrics, structured logging (`log/slog`), `/healthz` endpoint | — |
| 7 | Rate limiting: per-IP token bucket, `429` with `Retry-After` header | — |
| 8 | Forecast + compare: 5-day forecast endpoint, compare cities endpoint, Recharts frontend | — |
| 9 | Integration test suite: testcontainers (Redis + fake weather server), all failure scenarios | — |
| 10 | Live deployment: Fly.io, GitHub Actions deploy on merge to main | — |

---

## Step Detail

### Step 1 — Foundation
Go API with OpenWeatherMap adapter, in-memory cache, chi router, CORS, Swagger docs, mock-based unit tests, graceful shutdown on SIGTERM.

---

### Step 2 — React Frontend
Replaced Go HTML template with Vite + React + TypeScript. `WeatherCard` component replicates the old display. Go serves `web/dist/` in production; Vite dev server proxies `/api/*` to Go in development. Docker Compose split into two services: `api` (`:8080`) and `web` (`:5173`).

| File | Role |
|---|---|
| `web/src/api/weather.ts` | Typed fetch wrapper for `/api/weather` |
| `web/src/types/weather.ts` | TypeScript types mirroring Go DTOs |
| `web/src/components/WeatherCard.tsx` | Displays all weather fields |
| `web/vite.config.ts` | Proxy + `host: true` for container binding |
| `Dockerfile.api` | Go-only image for docker-compose |
| `web/Dockerfile` | Vite dev server image |

**Done when:** `make up` starts both services; weather card loads at `localhost:5173`.

---

### Step 3 — Docker + CI
Production `Dockerfile` already done (multi-stage: Node → Go → alpine). Remaining: GitHub Actions pipeline — test, lint, build on every push to main.

**Done when:** CI passes on push; `docker compose up` (via `make up`) starts both services cleanly.

---

### Step 4 — Multiple Providers + Circuit Breaker
Add Open-Meteo adapter (free, no API key). Provider router tries OWM first; circuit breaker opens on repeated failures and routes to Open-Meteo. Three states: closed → open → half-open → closed.

| Component | Responsibility |
|---|---|
| `internal/adapters/openmeteo.go` | Open-Meteo HTTP adapter |
| `internal/adapters/router.go` | Primary/fallback selection logic |
| `internal/adapters/breaker.go` | Circuit breaker state machine |

**Done when:** OWM key set to invalid → weather still returns from Open-Meteo. Logs show active provider.

**Unit tests:** breaker state transitions; router picks fallback on primary error.
**Integration tests:** fake HTTP server returns 500 for OWM; assert Open-Meteo adapter called.

---

### Step 5 — Redis Cache
Replace in-memory cache with Redis. Stale-while-revalidate: return stale data immediately, refresh in the background. `docker-compose.yml` adds Redis service (currently commented out).

**Done when:** `redis-cli` shows cached city keys; stale entry served on expiry, refreshed in background.

**Unit tests:** stale-while-revalidate — mock Redis returns expired entry, assert stale served + background refresh triggered.
**Integration tests** (testcontainers Redis): `Set` → `Get` round-trip; expired entry triggers background refresh.

---

### Step 6 — Observability
Prometheus metrics at `/metrics`: request count, duration histogram, cache hit/miss ratio, provider error count. Structured logging with `log/slog`. Health check at `/healthz` verifying cache and provider reachability.

**Done when:** `curl /metrics` returns Prometheus counters; `curl /healthz` returns `{"status":"ok"}`.

---

### Step 7 — Rate Limiting
Per-IP token bucket middleware (10 req/min). Returns `429 Too Many Requests` with `Retry-After` header.

**Done when:** 11th rapid request from same IP returns 429 with correct header.

**Unit tests:** middleware allows up to limit; blocks above it; `Retry-After` value is correct.

---

### Step 8 — Forecast + Compare Cities
New API endpoints and React pages:

| Endpoint | Description |
|---|---|
| `GET /api/forecast?city=<name>` | 5-day hourly forecast |
| `GET /api/compare?cities=<a>,<b>,<c>` | Parallel fetch for 2–4 cities |

React additions: forecast tab with Recharts temperature + rain line chart; compare page with side-by-side weather cards.

**Done when:** forecast chart shows 5 days of data; compare page shows cities side by side.

---

### Step 9 — Integration Test Suite
Full end-to-end tests with testcontainers (Redis + fake weather HTTP server):

| Scenario | Expected |
|---|---|
| Cache miss | Provider fetched, result cached |
| Stale cache | Stale served immediately, background refresh triggered |
| Provider failure | Circuit breaker opens, fallback provider used |
| Rate limit under load | Requests above limit return 429 |
| Health check with degraded provider | `/healthz` returns degraded status |

---

### Step 10 — Live Deployment
Fly.io deployment. GitHub Actions deploy step triggers on merge to main. `fly.toml` committed to repo.

**Done when:** `https://<app>.fly.dev/api/weather?city=london` returns live data; link in README.
