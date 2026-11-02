import type { WeatherResponse } from '../types/weather'

export async function fetchWeather(city: string): Promise<WeatherResponse> {
  const res = await fetch(`/api/weather?city=${encodeURIComponent(city)}`)
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.message ?? `Request failed: ${res.status}`)
  }
  return res.json()
}
