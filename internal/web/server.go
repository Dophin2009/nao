package web

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/Dophin2009/nao/pkg/api"
)

// HTTPReciever is a type alias for functions that handle
// HTTP requests.
type HTTPReciever = func(http.ResponseWriter, *http.Request, httprouter.Params)

// Handler is a single HTTP request handler
type Handler struct {
	Method          string
	Path            []string
	Func            HTTPReciever
	ResponseHeaders map[string]string
}

// PathString returns the full string form
// of the path of the handler
func (h *Handler) PathString() string {
	var str strings.Builder
	str.WriteString("/")
	for _, s := range h.Path {
		str.WriteString(s + "/")
	}
	return str.String()
}

// HandlerFunc returns a HTTP handler function that implements
// the handler's logic.
func (h *Handler) HandlerFunc() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		for k, v := range h.ResponseHeaders {
			w.Header().Add(k, v)
		}
		// Execute logic of handler
		if h.Func != nil {
			h.Func(w, r, ps)
		}
	}
}

// HandlerGroup is a group of handlers that have some
// shared properties
type HandlerGroup interface {
	Handlers() []Handler
}

const (
	// HeaderContentType is a HTTP header name that states
	// the structure of the response body
	HeaderContentType = "Content-Type"
	// HeaderContentTypeValJSON is a value for the content
	// type header for JSON
	HeaderContentTypeValJSON = "application/json"
)

// Server represents the API controller layer
type Server struct {
	Router  *httprouter.Router
	Address string
}

// NewServer returns a new instance of Controller
func NewServer(address string) Server {
	// Instantiate controller
	router := httprouter.New()
	s := Server{
		Router:  router,
		Address: address,
	}

	// Map routing handlers
	s.RegisterHandler(s.StatusHandler())

	return s
}

// HTTPServer returns a new http.Server object
// for this server
func (s *Server) HTTPServer() http.Server {
	return http.Server{
		Addr:    s.Address,
		Handler: s.Router,
	}
}

// RegisterHandler registers the given handler with the
// server
func (s *Server) RegisterHandler(h Handler) {
	s.Router.Handle(h.Method, h.PathString(), h.HandlerFunc())
}

// RegisterHandlerGroup registers all the handlers in the
// given handler group with the server
func (s *Server) RegisterHandlerGroup(g HandlerGroup) {
	for _, h := range g.Handlers() {
		s.RegisterHandler(h)
	}
}

// StatusHandler returns an endpoint handler that
// returns the status of the server
func (s *Server) StatusHandler() Handler {
	return Handler{
		Method: http.MethodGet,
		Path:   []string{},
		Func: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			status := api.StatusGet()
			json.NewEncoder(w).Encode(status)
		},
		ResponseHeaders: map[string]string{
			HeaderContentType: HeaderContentTypeValJSON,
		},
	}
}

// AuthenticationError is raised when the user
// fails to authenticate
type AuthenticationError struct {
	Debug string
}

func (err *AuthenticationError) Error() string {
	return fmt.Sprintf("error authenticating: %s", err.Debug)
}

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

const (
	// ErrorAuthentication is the generic error
	// message given when the user failed to
	// authenticate
	ErrorAuthentication = "error authenticating user"

	// ErrorPathVariableParsing is the generic
	// error message given when some path variable
	// could not be parsed properly
	ErrorPathVariableParsing = "error parsing path variable"

	// ErrorRequestBodyReading is the generic
	// error message given when HTTP request
	// body could not be read
	ErrorRequestBodyReading = "error reading request body"

	// ErrorRequestBodyParsing is the generic
	// error message given when HTTP request
	// body could not be parsed
	ErrorRequestBodyParsing = "error parsing request body"

	// ErrorInternalServer is the generic
	// error message given when an error was
	// encountered in the server
	ErrorInternalServer = "error within server"
)

// ReadRequestBody reads and returns the request body of the
// given HTTP request.
func ReadRequestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// ParsePathVar returns the string value of a path variable
// with the given name.
func ParsePathVar(varName string, ps *httprouter.Params) (value string, err error) {
	value = ps.ByName(varName)
	if value == "" {
		return value, errors.New("no such variable ''" + varName + "'")
	}
	return value, nil
}

// ParsePathVarInt returns the int value of a path variable
// with the given name.
func ParsePathVarInt(varName string, ps *httprouter.Params) (value int, err error) {
	v, err := ParsePathVar(varName, ps)
	if err != nil {
		return
	}

	value, err = strconv.Atoi(v)
	if err != nil {
		return
	}

	return
}

// EncodeResponseBody encodes the given value into the response
// body of the given ResponseWriter.
func EncodeResponseBody(body interface{}, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(body)
}

// EncodeResponseError encodes an error response into the response
// body of the given ResponseWriter.
func EncodeResponseError(err string, debug error, statusCode int, w http.ResponseWriter) {
	errorResponse := api.ErrorResponseNew(err, debug)
	json.NewEncoder(w).Encode(errorResponse)
	w.WriteHeader(statusCode)
}

// EncodeResponseErrorBadRequest encodes an error response with
// status code BadRequest.
func EncodeResponseErrorBadRequest(err string, debug error, w http.ResponseWriter) {
	EncodeResponseError(err, debug, http.StatusBadRequest, w)
}

// EncodeResponseErrorInternalServer encodes an error response
// with status code InternalServerError
func EncodeResponseErrorInternalServer(err string, debug error, w http.ResponseWriter) {
	EncodeResponseError(err, debug, http.StatusInternalServerError, w)
}

// EncodeResponseErrorUnauthorized encodes an error response with
// status code Unauthorized
func EncodeResponseErrorUnauthorized(err string, debug error, w http.ResponseWriter) {
	EncodeResponseError(err, debug, http.StatusUnauthorized, w)
}
