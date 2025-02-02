package adapters

import (
	"net/http"
	"weather-api/internal/apperrors"
	"weather-api/internal/config"
	"weather-api/internal/dto"

	"github.com/go-resty/resty/v2"
)

// openWeatherMapAdapter is an adapter for fetching third-party weather data
type openWeatherMapAdapter struct {
	client *resty.Client
}

// NewOpenWeatherMapAdapter creates a new openWeatherMapAdapter
func NewOpenWeatherMapAdapter() *openWeatherMapAdapter {
	return &openWeatherMapAdapter{
		client: resty.New(),
	}
}

// GetWeather fetches weather data from the OpenWeatherMap API
func (o *openWeatherMapAdapter) GetWeather(city string) (dto.OpenWeatherMapResponse, error) {
	weatherURL := config.GetConfig().OpenWeatherMapBaseURL + "/weather"

	var weatherData dto.OpenWeatherMapResponse
	resp, err := o.client.R().
		SetQueryParams(map[string]string{
			"q":     city,
			"appid": config.GetConfig().WeatherAPIKey,
			"units": "metric",
		}).
		SetResult(&weatherData).
		Get(weatherURL)

	if err != nil {
		return dto.OpenWeatherMapResponse{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return dto.OpenWeatherMapResponse{}, apperrors.ErrOpenWeatherMapFailure(resp.StatusCode())
	}

	return weatherData, nil
}
