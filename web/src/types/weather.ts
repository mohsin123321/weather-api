export interface WeatherResponse {
  coord: { lon: number; lat: number }
  weather: Array<{ id: number; main: string; description: string; icon: string }>
  main: {
    temp: number
    feels_like: number
    temp_min: number
    temp_max: number
    pressure: number
    humidity: number
  }
  wind: { speed: number; deg: number }
  clouds: { all: number }
  visibility: number
  rain?: { '1h': number }
  snow?: { '1h': number }
  sys: { sunrise: number; sunset: number; country: string }
  name: string
}
