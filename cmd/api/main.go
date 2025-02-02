package main

import (
	"log"
	"net/http"
	"weather-api/internal/adapters"
	"weather-api/internal/persistence"
	"weather-api/internal/server"

	"weather-api/internal/config"
)

func main() {
	// read the application configuration
	config.ReadConfig()

	cfg := config.GetConfig()

	server := server.NewServer(cfg.Port, persistence.NewCache(), adapters.NewOpenWeatherMapAdapter())

	// create a buffered channel used to notify the main goroutine that the shutdown is complete
	done := make(chan bool, 1)

	// go routine to gracefully shutdown the server
	go server.GracefulShutDown(done)

	// ignore http.ErrServerClosed since it indicates a graceful shutdown, which is expected behavior.
	if err := server.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server on port %s due to error: %v", cfg.Port, err)
	}

	// block until the graceful shutdown is complete
	<-done

	log.Println("Server shutdown gracefully")
}
