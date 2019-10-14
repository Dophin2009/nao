package api

import "time"

// Status contains information about the API
// at the current time
type Status struct {
	Version string     `json:"version"`
	Time    *time.Time `json:"time"`
}

// StatusGet retrieves information about the
// API at the current time and returns it as
// an APIStatus object
func StatusGet() *Status {
	currentTime := time.Now()
	return &Status{
		Version: "v1",
		Time:    &currentTime,
	}
}

// ErrorResponse represents an error message
// to be returned to the client if an error is
// encountered
type ErrorResponse struct {
	Time  *time.Time `json:"time"`
	Error string     `json:"error"`
	Debug string     `json:"debug"`
}

// ErrorResponseNew returns a new instance of
// errorResponse for the current time
func ErrorResponseNew(err string, debug error) *ErrorResponse {
	currentTime := time.Now()
	return &ErrorResponse{
		Time:  &currentTime,
		Error: err,
		Debug: debug.Error(),
	}
}
