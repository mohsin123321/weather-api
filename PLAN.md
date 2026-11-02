# Weather API — Implementation Plan

## Current State

A working Go weather API that proxies OpenWeatherMap, with in-memory caching, chi router, mock-based unit tests, and a vanilla JS + Go template frontend.

**Goal:** evolve it into a production-grade full-stack application that demonstrates real backend design (multiple providers, resilience, observability) and a modern React frontend.

---

## Architecture

```
Browser (React + Vite)
        │ HTTP (dev: :5173 → proxy → :8080)  (prod: Go serves dist/)
        ▼
  cmd/api (Go HTTP server)
        │
   ┌────┴────┐
   │  cache  │  in-memory → Redis (Step 5)
   └────┬────┘
        │
  provider router  ←── circuit breaker (Step 4)
   ┌────┴────┐
   │         │
 OWM      Open-Meteo   (Step 4)
```

---

## Directory Layout

```
/
├── cmd/api/              # entry point
├── internal/
│   ├── adapters/         # weather provider adapters
│   ├── apperrors/
│   ├── config/
│   ├── domain/           # interfaces
│   ├── dto/              # API response types
│   ├── persistence/      # cache (in-memory → Redis)
│   └── server/           # HTTP handlers, middleware
├── web/                  # React + TypeScript (Vite)
│   ├── src/
│   │   ├── api/          # typed fetch wrappers
│   │   ├── components/   # WeatherCard, ForecastChart, CityComparison
│   │   └── types/        # mirrors Go DTOs
│   └── dist/             # built output, served by Go in prod
├── tests/
├── Dockerfile
├── docker-compose.yml
└── .github/workflows/    # CI (Step 3)
```

---

## Build Steps

### Step 1 — Foundation (done)
Go API with OpenWeatherMap adapter, in-memory cache, chi router, CORS, Swagger, mock-based unit tests, graceful shutdown.

---

### Step 2 — React Frontend (current)
Replace Go HTML template with Vite + React + TypeScript. Replicate existing weather card as a React component. Go serves `web/dist/` as static files in production. Vite dev server proxies `/api/*` to Go in development.

**Done when:** `npm run dev` shows the weather card; `go run ./cmd/api/main.go` serves the built React app from `web/dist/`.

**Unit tests:** component renders weather data correctly (Vitest + React Testing Library).

---

### Step 3 — Docker + CI
`Dockerfile` (multi-stage: Node build → Go build → alpine runtime). `docker-compose.yml` for local development. GitHub Actions pipeline: test, lint, build on every push.

**Done when:** `docker compose up` starts the API; CI passes on push to main.

---

### Step 4 — Multiple Providers + Circuit Breaker
Add Open-Meteo adapter (free, no API key). Provider router tries primary (OWM) first; on failure the circuit breaker opens and routes to Open-Meteo. Circuit breaker has three states: closed, open, half-open.

**Done when:** with OWM key set to an invalid value, weather still returns (from Open-Meteo). Logs show which provider served the request.

**Unit tests:** circuit breaker state transitions (closed→open→half-open→closed). Provider router selects fallback when primary errors.

**Integration tests:** fake HTTP server returns 500 for OWM; assert response comes from Open-Meteo adapter.

---

### Step 5 — Redis Cache
Replace in-memory cache with Redis. Stale-while-revalidate: serve stale data immediately and refresh in the background. `docker-compose.yml` adds Redis service.

**Done when:** `CACHE_BACKEND=redis docker compose up` uses Redis; cached entries visible in `redis-cli`. Stale entry is served then refreshed in background.

**Unit tests:** stale-while-revalidate logic — mock Redis returns expired entry, assert stale served immediately + background refresh triggered.

**Integration tests** (testcontainers Redis): `Set` → `Get` round-trip; expired entry triggers background refresh.

---

### Step 6 — Observability
Prometheus metrics (`/metrics`): request count, request duration histogram, cache hit/miss ratio, provider error count. Structured logging with `log/slog`. Health check endpoint (`/healthz`): checks cache and provider reachability.

**Done when:** `curl /metrics` returns Prometheus-formatted counters; `curl /healthz` returns `{"status":"ok"}`.

---

### Step 7 — Rate Limiting
Per-IP rate limiting middleware (token bucket, 10 req/min). Returns `429 Too Many Requests` with `Retry-After` header.

**Done when:** 11 rapid requests from the same IP — first 10 succeed, 11th returns 429.

**Unit tests:** middleware allows up to limit, blocks above it. `Retry-After` header is set correctly.

---

### Step 8 — Forecast + Compare Cities
New API endpoints:
- `GET /api/forecast?city=<name>` — 5-day hourly forecast
- `GET /api/compare?cities=<a>,<b>,<c>` — parallel fetch for multiple cities

React frontend additions:
- Forecast tab with Recharts temperature + rain line chart
- Compare page with side-by-side weather cards

**Done when:** forecast chart shows 5 days of data; compare page shows 2–4 cities side by side.

---

### Step 9 — Integration Test Suite
Full end-to-end tests with testcontainers (Redis + fake weather server):
- Cache miss → provider fetch → cache set
- Stale cache → serve stale → background refresh
- Provider failure → circuit breaker opens → fallback provider
- Rate limit enforcement under concurrent load
- Health check reflects degraded provider

---

### Step 10 — Live Deployment
Fly.io deployment. GitHub Actions adds a deploy step on merge to main. `fly.toml` committed to repo. Live demo URL in README.

**Done when:** `https://<app>.fly.dev/api/weather?city=london` returns live data.
