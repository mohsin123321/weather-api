package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// responseWriter wraps the default http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader overrides the default WriteHeader method to capture the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Status returns the status code of the response
func (rw *responseWriter) Status() int {
	return rw.statusCode
}

// applyCommonMiddlewares apllies the common middlewares for the server.
func applyCommonMiddlewares(r *chi.Mux) {
	r.Use(corsMiddleware)
	r.Use(middleware.NoCache)
	r.Use(loggerMiddleware)
	r.Use(recoveryPanicMdlw)
}

// corsMiddleware set CORS for requests.
func corsMiddleware(next http.Handler) http.Handler {
	cors := cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodPut, http.MethodPatch, http.MethodOptions, http.MethodConnect, http.MethodTrace},
		AllowedHeaders:     []string{"*"},
		ExposedHeaders:     []string{"*"},
		AllowCredentials:   false,
		MaxAge:             0,
		OptionsPassthrough: false,
		Debug:              false,
	})

	return cors(next)
}

// loggerMiddleware log success and error responses details.
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// wrapper for the response writer to capture the status code
		rww := &responseWriter{w, http.StatusOK}

		defer func() {
			elapsedTime := time.Since(start).String()
			log.Println(r.Method, r.RemoteAddr, rww.Status(), r.URL.Path, elapsedTime)
		}()

		next.ServeHTTP(rww, r)
	})
}

// recover the panic called by an api
func recoveryPanicMdlw(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			err, ok := rec.(error)
			if !ok {
				// if the panic is not an error type then set err to unknown error
				err = errors.New("unknown error")
			}

			WriteErrorResponse(w, err)
		}()

		h.ServeHTTP(w, r)
	})
}
