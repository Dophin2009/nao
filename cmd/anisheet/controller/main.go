package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/Dophin2009/anisheet/pkg/api"
	bolt "go.etcd.io/bbolt"
)

// Controller represents the API controller layer
type Controller struct {
	DB     *bolt.DB
	Router *mux.Router
}

// NewController returns a new instance of Controller
func NewController(db *bolt.DB) Controller {
	router := mux.NewRouter().StrictSlash(true)

	c := Controller{
		DB:     db,
		Router: router,
	}

	c.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := api.StatusGet()
		json.NewEncoder(w).Encode(status)
	})

	mediaSubrouter := c.Router.PathPrefix("/media").Subrouter()
	mediaSubrouter.HandleFunc("/{id}", c.MediaQueryByID).Methods(http.MethodGet)
	mediaSubrouter.HandleFunc("/", c.MediaQueryAll).Methods(http.MethodGet)
	mediaSubrouter.HandleFunc("/", c.MediaCreate).Methods(http.MethodPost)

	return c
}

func encodeResponseBody(body interface{}, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(body)
}

func encodeError(err string, debug error, w http.ResponseWriter) {
	errorResponse := api.ErrorResponseNew(err, debug)
	json.NewEncoder(w).Encode(errorResponse)
	return
}

func withDefaultResponseHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Add("Content-Type", "application/json")
	return w
}
