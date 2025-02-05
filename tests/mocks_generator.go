package tests

// Adapters
//go:generate mockgen -package mock_adapters -destination=mocks/adapters/weather.go  weather-api/internal/domain OpenWeatherMapAdapter

// Persistence
//go:generate mockgen -package mock_persistence -destination=mocks/persistence/cache.go  weather-api/internal/domain CacheMap
