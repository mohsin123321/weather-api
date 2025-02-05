package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"text/template"
	"time"
	"weather-api/internal/apperrors"
	"weather-api/internal/domain"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Server wrapper of the http Server and also manages handlers which dispatches the requests to the appropriate routes.
type Server struct {
	httpServer *http.Server

	CacheMap              domain.CacheMap
	OpenWeatherMapAdapter domain.OpenWeatherMapAdapter
}

// NewServer is a factory to instantiate a new Server.
func NewServer(port string, cacheMap domain.CacheMap, OpenWeatherMapAdapter domain.OpenWeatherMapAdapter) *Server {
	newServer := &Server{
		CacheMap:              cacheMap,
		OpenWeatherMapAdapter: OpenWeatherMapAdapter,
	}

	newServer.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return newServer
}

// RegisterRoutes registers all HandlerFuncs for the existing HTTP routes.
func (s *Server) RegisterRoutes() http.Handler {
	router := chi.NewRouter()
	applyCommonMiddlewares(router)

	router.Get("/", http.HandlerFunc(s.HomePage))
	// Serve static assets
	router.Mount("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	router.Get("/api/docs/*", httpSwagger.WrapHandler)
	router.Get("/api/weather", http.HandlerFunc(s.GetWeatherDetails))

	return router
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	log.Printf("Starting server on port %s", s.httpServer.Addr[1:])
	return s.httpServer.ListenAndServe()
}

// GracefulShutDown listens for an interrupt signal from the OS and gracefully shuts down the server.
func (s *Server) GracefulShutDown(done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Block until a signal is received.
	<-ctx.Done()

	log.Println("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

// HomePage serves the home page by rendering an HTML template at the root path.
func (s *Server) HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.gohtml")
	if err != nil {
		WriteErrorResponse(w, apperrors.ErrServerError())
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		WriteErrorResponse(w, apperrors.ErrServerError())
		return
	}
}
