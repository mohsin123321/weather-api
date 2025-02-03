package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"weather-api/internal/apperrors"
)

// writeInternalError writes a raw HTTP 500 response.
func writeInternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	// manually write the response body in case the JSON marshalling fails in order to make the response consistent
	w.Write([]byte(`{"message": "Internal Server Error"}`))
	apperrors.ErrServerError().Log(apperrors.LogMessageUnexpectedError)
}

// WriteErrorResponse writes an app error as an HTTP response.
func WriteErrorResponse(w http.ResponseWriter, err error) {
	var appErr *apperrors.AppError
	errType := apperrors.LogMessageErrorResponse

	// If the error is not an app error, default to a server error
	if !errors.As(err, &appErr) {
		errType = apperrors.LogMessageUnexpectedError
		appErr = apperrors.ErrServerError()
	}

	bytes, err := appErr.MarhsalJSON()
	if err != nil {
		writeInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status())
	w.Write(bytes)

	appErr.Log(errType)
}

// WriteAPIDataResponse writes a successful API response with the given data.
func WriteAPIDataResponse(w http.ResponseWriter, code int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		writeInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}

// RecoverFromPanic recovers from a panic that occurs in the spawned goroutine.
func RecoverFromPanic(msg string) {
	if r := recover(); r != nil {
		log.Printf(msg, r)
	}
}
