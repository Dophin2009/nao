package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/Dophin2009/anisheet/pkg/api"
	"gitlab.com/Dophin2009/anisheet/pkg/data"
	bolt "go.etcd.io/bbolt"
)

// Controller represents the API controller layer
type Controller struct {
	Router       *mux.Router
	MediaService *data.MediaService
}

// New returns a new instance of Controller
func New(db *bolt.DB) Controller {
	// Instantiate controller
	router := mux.NewRouter().StrictSlash(true)
	c := Controller{
		Router: router,
		MediaService: &data.MediaService{
			DB: db,
		},
	}

	// Map routing handlers
	c.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := api.StatusGet()
		json.NewEncoder(w).Encode(status)
	})

	mediaSubrouter := c.Router.PathPrefix("/media").Subrouter()
	mediaSubrouter.HandleFunc("/", c.MediaCreate).Methods(http.MethodPost)
	mediaSubrouter.HandleFunc("/{id}", c.MediaQueryByID).Methods(http.MethodGet)
	mediaSubrouter.HandleFunc("/", c.MediaQueryAll).Methods(http.MethodGet)

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
