package server

import (
	"net/http"
	"strings"
	"weather-api/internal/apperrors"
	"weather-api/internal/config"
)

// @Summary Get weather details
// @ID GetWeatherDetails
// @tags Weather
// @Param city query string true "City name"
// @Success 200 {object} dto.OpenWeatherMapResponse
// @Produce json
// @Failure 400 "missing city query parameter"
// @Failure 500 "internal server error"
// @Router /weather [GET]
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

	// otherwise, delete the cache entry
	s.CacheMap.Delete(city)

	// Otherwise, fetch the weather from the API
	weather, err = s.OpenWeatherMapAdapter.GetWeather(city)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}

	// Cache the weather for future requests
	s.CacheMap.Set(city, weather, config.GetConfig().CacheExpiry)

	// Return the weather
	WriteAPIDataResponse(w, http.StatusOK, weather)
}
