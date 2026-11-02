import { useState } from 'react'
import { fetchWeather } from './api/weather'
import WeatherCard from './components/WeatherCard'
import type { WeatherResponse } from './types/weather'
import './App.css'

export default function App() {
  const [city, setCity] = useState('')
  const [weather, setWeather] = useState<WeatherResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!city.trim()) return

    setLoading(true)
    setError(null)

    try {
      const data = await fetchWeather(city.trim())
      setWeather(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
      setWeather(null)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="app">
      <h1 className="app-title">Weather</h1>

      <form className="search-form" onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Enter city name..."
          value={city}
          onChange={e => setCity(e.target.value)}
          disabled={loading}
        />
        <button type="submit" disabled={loading || !city.trim()}>
          {loading ? 'Loading...' : 'Search'}
        </button>
      </form>

      {error && <p className="error">{error}</p>}
      {weather && <WeatherCard data={weather} />}
    </div>
  )
}
