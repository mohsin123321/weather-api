package main

import (
	"log"
	"net/http"
	"time"
	"weather-api/internal/adapters"
	"weather-api/internal/persistence"
	"weather-api/internal/server"

	"weather-api/internal/config"
)

func main() {
	// read the application configuration
	config.ReadConfig()

	cfg := config.GetConfig()
	cache := persistence.NewCache()

	// start the cache cleanup routine in the background
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		// recover from a panic in the cache cleanup routine
		// prevents the panic from crashing the application
		defer server.RecoverFromPanic("recovered from panic in cache cleanup routine: %v")
		// run the cleanup routine after every minute
		cache.CleanUp(ticker)
	}()

	s := server.NewServer(cfg.Port, cache, adapters.NewOpenWeatherMapAdapter())

	// create a buffered channel used to notify the main goroutine that the shutdown is complete
	done := make(chan bool, 1)

	// go routine to gracefully shutdown the server
	go s.GracefulShutDown(done)

	// ignore http.ErrServerClosed since it indicates a graceful shutdown, which is expected behavior.
	if err := s.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server on port %s due to error: %v", cfg.Port, err)
	}

	// block until the graceful shutdown is complete
	<-done

	log.Println("Server shutdown gracefully")
}
