package fakeddto

import "weather-api/internal/dto"

// OpenWeatherMapResponse returns a fake OpenWeatherMapResponse
func OpenWeatherMapResponse() dto.OpenWeatherMapResponse {
	rain := dto.Rain{OneHour: 0.25}
	snow := dto.Snow{OneHour: 0.25}
	return dto.OpenWeatherMapResponse{
		ID:    2643743,
		Coord: dto.Coord{Lon: -0.13, Lat: 51.51},
		Weather: []dto.Weather{
			{ID: 300, Main: "Drizzle", Description: "light intensity drizzle", Icon: "09d"},
		},
		Base:       "stations",
		Main:       dto.Main{Temp: 10, FeelsLike: 8.5, Pressure: 1012, Humidity: 81, TempMin: 10, TempMax: 10},
		Visibility: 10000,
		Wind:       dto.Wind{Speed: 4.1, Deg: 80, Gust: 8.2},
		Clouds:     dto.Clouds{All: 90},
		Dt:         1485789600,
		Sys:        dto.Sys{Type: 1, ID: 5091, Country: "GB", Sunrise: 1485762037, Sunset: 1485794875},
		Rain:       &rain,
		Snow:       &snow,
		Timezone:   0,
		Name:       "London",
		Cod:        200,
	}
}
