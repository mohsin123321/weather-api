package domain

import "weather-api/internal/dto"

// OpenWeatherMapAdapter is an interface for getting weather data
type OpenWeatherMapAdapter interface {
	GetWeather(city string) (dto.OpenWeatherMapResponse, error)
}

// CacheMap is an interface for caching weather data
type CacheMap interface {
	Get(city string) (dto.OpenWeatherMapResponse, bool, error)
	Set(city string, weather dto.OpenWeatherMapResponse, cacheExpiry int)
	Delete(city string)
}
