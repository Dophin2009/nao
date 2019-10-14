package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/Dophin2009/anisheet/pkg/data"
)

// MediaQueryByID gets a single Media by ID from the
// persistence layer and writes it to the HTTP response
func (c *Controller) MediaQueryByID(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	var vars map[string]string = mux.Vars(r)
	idVal := vars["id"]

	id, err := strconv.Atoi(idVal)
	if err != nil {
		encodeError("error parsing id"+idVal, err, w)
		return
	}

	media, err := data.MediaGet(id, c.DB)
	if err != nil {
		encodeError("error querying media", err, w)
		return
	}

	json.NewEncoder(w).Encode(media)
}

// MediaQueryAll gets all the Media from the
// persistence layer and writes it to the HTTP response
func (c *Controller) MediaQueryAll(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	media, err := data.MediaGetAll(c.DB)
	if err != nil {
		encodeError("error querying media", err, w)
		return
	}
	if media == nil {
		media = []data.Media{}
	}

	json.NewEncoder(w).Encode(media)
}

// MediaCreate persists the request body
func (c *Controller) MediaCreate(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encodeError("error reading request body", err, w)
		return
	}

	media := data.Media{}
	err = json.Unmarshal(body, &media)
	if err != nil {
		encodeError("error parsing request body", err, w)
		return
	}

	err = data.MediaCreate(&media, c.DB)
	if err != nil {
		encodeError("error creating media", err, w)
		return
	}

	json.NewEncoder(w).Encode(media)
}
