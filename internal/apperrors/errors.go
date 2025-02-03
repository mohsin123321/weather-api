package apperrors

import (
	"encoding/json"
	"log"
)

// Log message types are used to define the type of each log.
const (
	LogMessageErrorResponse   = "ERROR RESPONSE"   // Handled errors that could happen.
	LogMessageUnexpectedError = "UNEXPECTED ERROR" // Panic errors that should never happen.
)

// AppError represents an application error.
type AppError struct {
	msg    string
	status int
}

// Error returns the error message.
func (e *AppError) Error() string {
	return e.msg
}

// Status returns the HTTP status code.
func (e *AppError) Status() int {
	return e.status
}

// MarhsalJSON marshals the error message into a JSON byte slice.
func (e *AppError) MarhsalJSON() ([]byte, error) {
	json, error := json.Marshal(struct{ Message string }{Message: e.msg})
	return json, error
}

// Log logs the error message.
func (e *AppError) Log(msg string) {
	log.Printf("%s: %s %d ", msg, e.Error(), e.Status())
}
