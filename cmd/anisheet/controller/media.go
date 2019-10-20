package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/Dophin2009/anisheet/pkg/api"
	"gitlab.com/Dophin2009/anisheet/pkg/data"
)

// MediaQueryByID gets a single Media by ID (given by the path
// variable {id}) from the persistence layer and writes it
// to the HTTP response
func (c *Controller) MediaQueryByID(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	var vars map[string]string = mux.Vars(r)
	idVal := vars["id"]

	id, err := strconv.Atoi(idVal)
	if err != nil {
		encodeError(api.PathVariableParsingError, err, w)
		return
	}

	var media = data.Media{
		ID: id,
	}
	err = c.MediaService.GetByID(&media)
	if err != nil {
		encodeError(api.DatabaseQueryingError, err, w)
		return
	}

	encodeResponseBody(media, w)
}

// MediaQueryAll gets all the Media from the
// persistence layer and writes it to the HTTP response
func (c *Controller) MediaQueryAll(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	media, err := c.MediaService.GetAll()
	if err != nil {
		encodeError(api.DatabaseQueryingError, err, w)
		return
	}
	if media == nil {
		media = []data.Media{}
	}

	encodeResponseBody(media, w)
}

// MediaCreate persists the request body
func (c *Controller) MediaCreate(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encodeError(api.RequestBodyReadingError, err, w)
		return
	}

	media := data.Media{}
	err = json.Unmarshal(body, &media)
	if err != nil {
		encodeError(api.RequestBodyParsingError, err, w)
		return
	}

	err = c.MediaService.Create(&media)
	if err != nil {
		encodeError(api.DatabasePersistingError, err, w)
		return
	}

	encodeResponseBody(media, w)
}

// MediaUpdate persists the request body to
// an existing ID
func (c *Controller) MediaUpdate(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encodeError(api.RequestBodyReadingError, err, w)
		return
	}

	media := data.Media{}
	err = json.Unmarshal(body, &media)
	if err != nil {
		encodeError(api.RequestBodyParsingError, err, w)
		return
	}

	err = c.MediaService.Update(&media)
	if err != nil {
		encodeError(api.DatabasePersistingError, err, w)
		return
	}

	encodeResponseBody(media, w)
}
