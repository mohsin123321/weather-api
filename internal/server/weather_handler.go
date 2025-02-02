package server

import (
	"net/http"
	"strings"
	"weather-api/internal/apperrors"
)

// GetWeatherDetails fetches weather data for a given city.
func (s *Server) GetWeatherDetails(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")

	if city == "" {
		WriteErrorResponse(w, apperrors.ErrMissingCityParam())
		return
	}

	// used string.ToLower to ensure map keys are case-insensitive
	city = strings.ToLower(city)
	weather, found, err := s.CacheMap.Get(city)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}

	// Check if we have the weather cached
	// If we do, return the cached weather
	if found {
		WriteAPIDataResponse(w, http.StatusOK, weather)
		return
	}

	// Otherwise, fetch the weather from the API
	weather, err = s.OpenWeatherMapAdapter.GetWeather(city)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}

	// Cache the weather for future requests
	s.CacheMap.Set(city, weather)

	// Return the weather
	WriteAPIDataResponse(w, http.StatusOK, weather)
}
