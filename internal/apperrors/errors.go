package apperrors

import "encoding/json"

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
