package apperrors

import (
	"net/http"
)

// ErrServerError is raised when server breaks for internal reasons.
func ErrServerError() *AppError {
	return &AppError{
		msg:    "Internal Server Error",
		status: http.StatusInternalServerError,
	}
}

// ErrMissingCityParam is raised when the city query parameter is missing from the request.
func ErrMissingCityParam() *AppError {
	return &AppError{
		msg:    "missing city query parameter",
		status: http.StatusBadRequest,
	}
}

// ErrOpenWeatherMapFailure is raised when the OpenWeatherMap API fails to fetch weather data.
func ErrOpenWeatherMapFailure(status int) *AppError {
	return &AppError{
		msg:    "failed to fetch weather data",
		status: status,
	}
}

// ErrInvalidCacheEntryType is raised when the cache entry type is invalid.
func ErrInvalidCacheEntryType() *AppError {
	return &AppError{
		msg:    "invalid cache entry type",
		status: http.StatusInternalServerError,
	}
}
