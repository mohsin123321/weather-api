package persistence_test

import (
	"os"
	"testing"
	"time"
	"weather-api/internal/config"
	"weather-api/internal/persistence"
	fakeddto "weather-api/tests/fakedto"
)

func TestGet(t *testing.T) {
	cache := persistence.NewCache()
	londonWeather := fakeddto.OpenWeatherMapResponse()
	londonWeather.Sys.Country = "GB"
	italyWeather := fakeddto.OpenWeatherMapResponse()
	italyWeather.Sys.Country = "IT"

	cache.Set("London", fakeddto.OpenWeatherMapResponse(), 5)
	cache.Set("Italy", fakeddto.OpenWeatherMapResponse(), 5)

	t.Run("Get key", func(t *testing.T) {
		weather, found, err := cache.Get("London")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		if !found {
			t.Errorf("Expected true, got %v", found)
		}
		if weather.Sys.Country != "GB" {
			t.Errorf("Expected GB, got %s", weather.Sys.Country)
		}
	})

	t.Run("Get key not found", func(t *testing.T) {
		weather, found, err := cache.Get("Paris")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		if found {
			t.Errorf("Expected false, got %v", found)
		}
		if weather.ID != 0 {
			t.Errorf("Expected 0, got %d", weather.ID)
		}
	})

	t.Run("Get Key Expired Entry", func(t *testing.T) {
		cfg := config.GetConfig()
		cfg.CacheExpiry = 0
		cache := persistence.NewCache()
		cache.Set("London", fakeddto.OpenWeatherMapResponse(), 0)

		_, found, err := cache.Get("London")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		if found {
			t.Errorf("Expected false, got %v", found)
		}
	})
}

func TestSet(t *testing.T) {
	t.Run("Set key", func(t *testing.T) {
		cache := persistence.NewCache()
		cache.Set("London", fakeddto.OpenWeatherMapResponse(), 5)
		weather, found, err := cache.Get("London")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		if !found {
			t.Errorf("Expected true, got %v", found)
		}
		if weather.ID != fakeddto.OpenWeatherMapResponse().ID {
			t.Errorf("Expected 1, got %d", weather.ID)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete key", func(t *testing.T) {
		cache := persistence.NewCache()
		cache.Set("London", fakeddto.OpenWeatherMapResponse(), 5)
		cache.Delete("London")
		_, found, _ := cache.Get("London")
		if found {
			t.Errorf("Expected false, got %v", found)
		}
	})
}

func TestCleanUp(t *testing.T) {
	t.Run("CleanUp", func(t *testing.T) {
		cache := persistence.NewCache()
		cache.Set("London", fakeddto.OpenWeatherMapResponse(), 5)
		cache.Set("Italy", fakeddto.OpenWeatherMapResponse(), 0)
		cache.Set("Paris", fakeddto.OpenWeatherMapResponse(), 0)
		ticker := time.NewTicker(5 * time.Millisecond)

		// run the cleanup routine in a separate goroutine so that the test can continue
		go cache.CleanUp(ticker)

		// wait for the cleanup to process
		time.Sleep(20 * time.Millisecond)

		_, found, _ := cache.Get("London")
		if !found {
			t.Errorf("Expected false, got %v", found)
		}

		_, found, _ = cache.Get("Italy")
		if found {
			t.Errorf("Expected false, got %v", found)
		}

		_, found, _ = cache.Get("Paris")
		if found {
			t.Errorf("Expected false, got %v", found)
		}
	})
}

func TestMain(m *testing.M) {
	// set environment variables for testing with default values
	config.ReadConfig()

	// Run the tests and exit with the result code
	os.Exit(m.Run())
}
