package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/Dophin2009/nao/pkg/api"
)

// Handler is a single HTTP request handler
type Handler struct {
	Method          string
	Path            []string
	Logic           func(http.ResponseWriter, *http.Request, httprouter.Params)
	Authentication  func(http.ResponseWriter, *http.Request, httprouter.Params) error
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

// AuthenticationError is raised when the user
// fails to authenticate
type AuthenticationError struct {
	Debug string
}

func (err *AuthenticationError) Error() string {
	return fmt.Sprintf("error authenticating user: %s", err.Debug)
}

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
	s.Router.Handle(h.Method, h.PathString(), func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Check authentication
		if h.Authentication != nil {
			if err := h.Authentication(w, r, ps); err != nil {
				// If error in authenticating, not any other reason,
				// return unauthorized status
				if err, ok := err.(*AuthenticationError); ok {
					encodeError("authentication failed", err, http.StatusUnauthorized, w)
					return
				}
				// Else, return internal error code
				encodeError("internal server error", err, http.StatusInternalServerError, w)
				return
			}
		}
		for k, v := range h.ResponseHeaders {
			w.Header().Add(k, v)
		}
		// Execute logic of handler
		if h.Logic != nil {
			h.Logic(w, r, ps)
		}
	})
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
		Logic: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			status := api.StatusGet()
			json.NewEncoder(w).Encode(status)
		},
		ResponseHeaders: map[string]string{
			HeaderContentType: HeaderContentTypeValJSON,
		},
	}
}

func parsePathVar(ps *httprouter.Params, name string) (value string, err error) {
	value = ps.ByName(name)
	if value == "" {
		return value, errors.New("no such variable ''" + name + "'")
	}
	return value, nil
}

func parsePathVarInteger(ps *httprouter.Params, name string) (value int, err error) {
	v, err := parsePathVar(ps, name)
	if err != nil {
		return
	}

	value, err = strconv.Atoi(v)
	if err != nil {
		return
	}

	return
}

func encodeResponseBody(body interface{}, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(body)
}

func encodeError(err string, debug error, statusCode int, w http.ResponseWriter) {
	errorResponse := api.ErrorResponseNew(err, debug)
	json.NewEncoder(w).Encode(errorResponse)
	w.WriteHeader(statusCode)
	return
}
