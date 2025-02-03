package config

import (
	"log"
	"strconv"
	"sync"

	"os"

	"github.com/joho/godotenv"
)

var (
	config Config
	once   sync.Once
)

// Config represents the application configuration
type Config struct {
	Port                  string
	WeatherAPIKey         string
	CacheExpiry           int
	OpenWeatherMapBaseURL string
}

// Get returns the application configuration
func GetConfig() Config {
	return config
}

// ReadConfig reads the application configuration from the environment
func ReadConfig() {
	// follow the singleton pattern to ensure the configuration is only loaded once
	once.Do(func() {
		loadConfig()
	})
}

// loadConfig loads the application configuration from the environment
func loadConfig() {
	// load the environment variables from .env file
	// Note: the environment variables take precedence over the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default values or system environment variables")
	}

	// initialize the configuration with env/default values
	config = Config{
		Port:                  getEnv("PORT", "8080"),
		WeatherAPIKey:         getEnv("WEATHER_API_KEY", ""),
		CacheExpiry:           getEnvInt("CACHE_EXPIRY", 5),
		OpenWeatherMapBaseURL: getEnv("OPEN_WEATHER_MAP_BASE_URL", "https://api.openweathermap.org/data/2.5"),
	}
}

// getEnv retrieves the value of an environment variable, or a fallback value if it's not set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvInt retrieves the value of an environment variable as an integer, or a fallback value if it's not set
func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}
