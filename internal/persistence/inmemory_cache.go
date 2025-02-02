package persistence

import (
	"sync"
	"time"
	"weather-api/internal/apperrors"
	"weather-api/internal/config"
	"weather-api/internal/dto"
)

// CacheMap is an in-memory cache for storing weather data
type CacheMap struct {
	data sync.Map
}

// cacheEntry is a struct to store weather data and its expiration time
type cacheEntry struct {
	weather   dto.OpenWeatherMapResponse
	expiresAt time.Time
}

// NewCache creates a new CacheMap
func NewCache() *CacheMap {
	return &CacheMap{
		data: sync.Map{},
	}
}

// Get retrieves weather data from the cache if it exists and is not expired
func (c *CacheMap) Get(city string) (dto.OpenWeatherMapResponse, bool, error) {
	cached, found := c.data.Load(city)
	if !found {
		return dto.OpenWeatherMapResponse{}, false, nil
	}

	entry, ok := cached.(cacheEntry)
	if !ok {
		return dto.OpenWeatherMapResponse{}, false, apperrors.ErrInvalidCacheEntryType()
	}

	if time.Now().After(entry.expiresAt) {
		c.data.Delete(city)
		return dto.OpenWeatherMapResponse{}, false, nil
	}

	return entry.weather, true, nil
}

// Set stores weather data in the cache with a 5-minute expiration
func (c *CacheMap) Set(city string, weather dto.OpenWeatherMapResponse) {
	// cacheExpiry in minutes loaded from the config
	cacheExpiry := config.GetConfig().CacheExpiry
	c.data.Store(city, cacheEntry{
		weather:   weather,
		expiresAt: time.Now().Add(time.Duration(cacheExpiry) * time.Minute),
	})
}
