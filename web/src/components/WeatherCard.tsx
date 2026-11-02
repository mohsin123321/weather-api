import type { WeatherResponse } from '../types/weather'

interface Props {
  data: WeatherResponse
}

function formatTime(unix: number): string {
  return new Date(unix * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

export default function WeatherCard({ data }: Props) {
  const main = data.weather[0]

  return (
    <div className="weather-card">
      <h1>{data.name}, {data.sys.country}</h1>
      <p className="coords">
        {data.coord.lat}°N {data.coord.lon}°E
      </p>

      <div className="condition">
        <img
          src={`https://openweathermap.org/img/wn/${main.icon}@2x.png`}
          alt={main.description}
        />
        <h2>{main.description}</h2>
      </div>

      <p className="temp">{Math.round(data.main.temp)}°C</p>
      <p className="feels-like">Feels like {Math.round(data.main.feels_like)}°C</p>

      <div className="details">
        <Detail label="High / Low" value={`${Math.round(data.main.temp_max)}° / ${Math.round(data.main.temp_min)}°`} />
        <Detail label="Humidity" value={`${data.main.humidity}%`} />
        <Detail label="Pressure" value={`${data.main.pressure} hPa`} />
        <Detail label="Wind" value={`${data.wind.speed} m/s`} />
        <Detail label="Clouds" value={`${data.clouds.all}%`} />
        <Detail label="Visibility" value={`${(data.visibility / 1000).toFixed(1)} km`} />
        {data.rain?.['1h'] !== undefined && (
          <Detail label="Rain (1h)" value={`${data.rain['1h']} mm`} />
        )}
        {data.snow?.['1h'] !== undefined && (
          <Detail label="Snow (1h)" value={`${data.snow['1h']} mm`} />
        )}
        <Detail label="Sunrise" value={formatTime(data.sys.sunrise)} />
        <Detail label="Sunset" value={formatTime(data.sys.sunset)} />
      </div>
    </div>
  )
}

function Detail({ label, value }: { label: string; value: string }) {
  return (
    <p className="detail-row">
      <span className="detail-label">{label}</span>
      <span className="detail-value">{value}</span>
    </p>
  )
}
