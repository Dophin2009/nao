package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cheekybits/genny/generic"
	"github.com/gorilla/mux"
	"gitlab.com/Dophin2009/nao/pkg/api"
	"gitlab.com/Dophin2009/nao/pkg/data"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "EntityType=Media,Episode,Character,Genre,Producer,Person,User,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia"

// EntityType is a generic placeholder for all entity types;
// it is assumed that EntityType structs have an ID and Version,
// both of which should be of int type.
type EntityType generic.Type

// EntityTypeSubrouter adds a subrouter with /EntityType as the
// path prefix to the given Controller
func EntityTypeSubrouter(c *Controller) *mux.Router {
	subrouter := c.Router.PathPrefix("/EntityType").Subrouter()
	subrouter.HandleFunc("", c.EntityTypeCreate).Methods(http.MethodPost)
	subrouter.HandleFunc("", c.EntityTypeUpdate).Methods(http.MethodPut)
	subrouter.HandleFunc("/{id}", c.EntityTypeDelete).Methods(http.MethodDelete)
	subrouter.HandleFunc("/{id}", c.EntityTypeQueryByID).Methods(http.MethodGet)
	subrouter.HandleFunc("", c.EntityTypeQueryAll).Methods(http.MethodGet)
	return subrouter
}

// EntityTypeQueryByID gets a single EntityType by ID (given
// by the path variable {id}) from the persistence layer and
// writes it to the HTTP response.
func (c *Controller) EntityTypeQueryByID(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	// Parse ID from path
	id, err := parseID(r)
	if err != nil {
		encodeError(api.PathVariableParsingError, err, w)
		return
	}

	// Retrieve EntityType by ID
	e := data.EntityType{
		ID: id,
	}
	err = c.EntityTypeService.GetByID(&e)
	if err != nil {
		encodeError(api.DatabaseQueryingError, err, w)
		return
	}

	// Encode response
	encodeResponseBody(e, w)
}

// EntityTypeQueryAll gets all the EntityType from the
// persistence layer and writes it to the HTTP response.
func (c *Controller) EntityTypeQueryAll(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	// Retrieve all EntityType by ID
	list, err := c.EntityTypeService.GetAll()
	if err != nil {
		encodeError(api.DatabaseQueryingError, err, w)
		return
	}
	if list == nil {
		list = []data.EntityType{}
	}

	// Encode response
	encodeResponseBody(list, w)
}

// EntityTypeCreate parses and persists the request body.
func (c *Controller) EntityTypeCreate(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encodeError(api.RequestBodyReadingError, err, w)
		return
	}

	// Parse request body into EntityType
	var e data.EntityType
	err = json.Unmarshal(body, &e)
	if err != nil {
		encodeError(api.RequestBodyParsingError, err, w)
		return
	}

	// Persist parsed EntityType
	err = c.EntityTypeService.Create(&e)
	if err != nil {
		encodeError(api.DatabasePersistingError, err, w)
		return
	}

	// Encode response
	encodeResponseBody(e, w)
}

// EntityTypeUpdate parses and persists the request body
// to an existing ID, replacing the existing EntityType.
func (c *Controller) EntityTypeUpdate(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encodeError(api.RequestBodyReadingError, err, w)
		return
	}

	// Parse request body into EntityType
	var e data.EntityType
	err = json.Unmarshal(body, &e)
	if err != nil {
		encodeError(api.RequestBodyParsingError, err, w)
		return
	}

	// Persist parsed EntityType
	err = c.EntityTypeService.Update(&e)
	if err != nil {
		encodeError(api.DatabasePersistingError, err, w)
		return
	}

	// Encode response
	encodeResponseBody(e, w)
}

// EntityTypeDelete deletes the EntityType of the given ID in the
// persistence layer and returns the existing value.
func (c *Controller) EntityTypeDelete(w http.ResponseWriter, r *http.Request) {
	w = withDefaultResponseHeaders(w)

	// Parse ID
	id, err := parseID(r)
	if err != nil {
		encodeError(api.PathVariableParsingError, err, w)
		return
	}

	// Delete by ID and retrieve existing
	e, err := c.EntityTypeService.Delete(id)
	if err != nil {
		encodeError(api.DatabasePersistingError, err, w)
		return
	}

	// Encode response
	encodeResponseBody(e, w)
}
