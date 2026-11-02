# Weather App

A weather application that fetches real-time weather data for any city using the OpenWeatherMap API. Built with a Go backend and a React + TypeScript frontend.

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.22, Chi router, in-memory cache |
| Frontend | React 18, TypeScript, Vite |
| Containerisation | Docker, Docker Compose |
| API | OpenWeatherMap |

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- An [OpenWeatherMap API key](https://openweathermap.org/api) (free tier is enough)

## Getting Started

**1. Copy the environment file and add your API key:**

```bash
cp .env.example .env
# then edit .env and set WEATHER_API_KEY=your_key_here
```

**2. Start all services:**

```bash
make up
```

This builds and starts the Go API (`localhost:8080`) and the React dev server (`localhost:5173`) in the background. Open `http://localhost:5173` in your browser.

**3. Stop all services:**

```bash
make down
```

**4. Tail live logs:**

```bash
make logs
```

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `WEATHER_API_KEY` | OpenWeatherMap API key (required) | — |
| `PORT` | Port the Go API listens on | `8080` |
| `CACHE_EXPIRY` | Cache TTL in minutes | `5` |
| `OPEN_WEATHER_MAP_BASE_URL` | OpenWeatherMap base URL | `https://api.openweathermap.org/data/2.5` |

## Running Locally (without Docker)

```bash
# Terminal 1 — Go API
go run ./cmd/api/main.go

# Terminal 2 — React dev server (hot reload)
cd web && npm run dev
```

## All Make Commands

| Command | Description |
|---|---|
| `make up` | Build and start all services via Docker Compose |
| `make down` | Stop and remove all containers |
| `make logs` | Follow live logs from all services |
| `make run` | Run the Go API locally (no Docker) |
| `make build` | Compile the Go binary to `bin/api` |
| `make test` | Generate mocks and run all tests |
| `make lint` | Run golangci-lint |
| `make docs` | Regenerate Swagger docs |

## Running Tests

Install `mockgen` once:

```bash
go install go.uber.org/mock/mockgen@latest
```

Then:

```bash
make test
```

## Build Plan

See [PLAN.md](PLAN.md) for the full roadmap.
