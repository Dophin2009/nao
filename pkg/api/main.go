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

// AuthenticationError is the generic error
// message given when the user failed to
// authenticate
const AuthenticationError = "error authenticating user"

// PathVariableParsingError is the generic
// error message given when some path variable
// could not be parsed properly
const PathVariableParsingError = "error parsing path variable"

// RequestBodyReadingError is the generic
// error message given when HTTP request
// body could not be read
const RequestBodyReadingError = "error reading request body"

// RequestBodyParsingError is the generic
// error message given when HTTP request
// body could not be parsed
const RequestBodyParsingError = "error parsing request body"

// DatabaseQueryingError is the generic
// error message given when an error was
// encountered while querying the database
const DatabaseQueryingError = "error querying database"

// DatabasePersistingError is the generic
// error message given when an error was
// encountered while persisting the database
const DatabasePersistingError = "error persisting in database"

// GenericInternalServerError is the generic
// error message given when an error was
// encountered in the server
const GenericInternalServerError = "internal server error"
