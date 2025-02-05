package server_test

import (
	"net/http"
	"strings"
	"testing"
	"weather-api/internal/apperrors"
	"weather-api/internal/config"
	"weather-api/internal/dto"
	fakeddto "weather-api/tests/fakedto"
)

func TestGetWeatherDetails(t *testing.T) {
	tester := setup(t)
	t.Run("Error city param missing", func(t *testing.T) {
		url := "/api/weather"
		req, w := mockReqGenerator(url)

		tester.Server.GetWeatherDetails(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", resp.StatusCode)
		}

		var errorMessage errorMessage
		unmarshalResponseBody(&errorMessage, resp, t)

		expectedMessage := "missing city query parameter"
		if errorMessage.Message != expectedMessage {
			t.Errorf("Expected message %s, got %s", expectedMessage, errorMessage.Message)
		}
	})

	t.Run("Error cacheMap.Get", func(t *testing.T) {
		city := "London"
		url := "/api/weather?city=" + city
		req, w := mockReqGenerator(url)

		city = strings.ToLower(city)
		tester.CacheMap.EXPECT().Get(city).Return(dto.OpenWeatherMapResponse{}, false, apperrors.ErrServerError())

		tester.Server.GetWeatherDetails(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status code 400, got %d", resp.StatusCode)
		}

		var errorMessage errorMessage
		unmarshalResponseBody(&errorMessage, resp, t)

		expectedMessage := "internal server error"
		if errorMessage.Message != expectedMessage {
			t.Errorf("Expected message %s, got %s", expectedMessage, errorMessage.Message)
		}
	})

	t.Run("Cache hit", func(t *testing.T) {
		city := "London"
		url := "/api/weather?city=" + city
		req, w := mockReqGenerator(url)

		city = strings.ToLower(city)
		tester.CacheMap.EXPECT().Get(city).Return(fakeddto.OpenWeatherMapResponse(), true, nil)

		tester.Server.GetWeatherDetails(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", resp.StatusCode)
		}

		var weatherDetails dto.OpenWeatherMapResponse
		unmarshalResponseBody(&weatherDetails, resp, t)

		expectedWeatherDetails := fakeddto.OpenWeatherMapResponse()

		if expectedWeatherDetails.ID != weatherDetails.ID || expectedWeatherDetails.Main.Temp != weatherDetails.Main.Temp {
			t.Errorf("Expected weather details %v, got %v", expectedWeatherDetails, weatherDetails)
		}
	})

	t.Run("Cache miss", func(t *testing.T) {
		city := "London"
		url := "/api/weather?city=" + city
		req, w := mockReqGenerator(url)

		city = strings.ToLower(city)
		tester.CacheMap.EXPECT().Get(city).Return(dto.OpenWeatherMapResponse{}, false, nil)
		tester.CacheMap.EXPECT().Delete(city)
		tester.OpenWeatherMapAdapter.EXPECT().GetWeather(city).Return(fakeddto.OpenWeatherMapResponse(), nil)
		tester.CacheMap.EXPECT().Set(city, fakeddto.OpenWeatherMapResponse(), config.GetConfig().CacheExpiry)

		tester.Server.GetWeatherDetails(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", resp.StatusCode)
		}

		var weatherDetails dto.OpenWeatherMapResponse
		unmarshalResponseBody(&weatherDetails, resp, t)

		expectedWeatherDetails := fakeddto.OpenWeatherMapResponse()

		if expectedWeatherDetails.ID != weatherDetails.ID || expectedWeatherDetails.Main.Temp != weatherDetails.Main.Temp {
			t.Errorf("Expected weather details %v, got %v", expectedWeatherDetails, weatherDetails)
		}
	})
}
