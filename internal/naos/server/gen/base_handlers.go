package server

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cheekybits/genny/generic"
	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/Dophin2009/nao/pkg/api"
	"gitlab.com/Dophin2009/nao/pkg/data"
)

//go:generate genny -in=$GOFILE -out=gen.$GOFILE gen "EntityType=Media,Episode,Character,Genre,Producer,Person,User,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia,UserMediaList"

// EntityType is a generic placeholder for all entity types;
// it is assumed that EntityType structs have an ID and Version,
// both of which should be of int type.
type EntityType generic.Type

// EntityTypeHandlerGroup is a basic handler group for EntityType
type EntityTypeHandlerGroup struct {
	Service *data.EntityTypeService
}

// NewEntityTypeHandlerGroup returns a handler group for
// EntityType with the given service
func NewEntityTypeHandlerGroup(service *data.EntityTypeService) EntityTypeHandlerGroup {
	g := EntityTypeHandlerGroup{
		Service: service,
	}
	return g
}

// Handlers returns all the basic CRUD handlers for the
// handler group
func (g *EntityTypeHandlerGroup) Handlers() []Handler {
	return []Handler{
		g.CreateHandler(),
		g.UpdateHandler(),
		g.DeleteHandler(),
		g.GetAllHandler(),
		g.GetByIDHandler(),
	}
}

// CreateHandler returns an POST endpoint handler
// for creating new EntityType
func (g *EntityTypeHandlerGroup) CreateHandler() Handler {
	return Handler{
		Method: http.MethodPost,
		Path:   []string{strings.ToLower("EntityType")},
		Logic: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Read request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				encodeError(api.RequestBodyReadingError, err, http.StatusBadRequest, w)
				return
			}

			// Parse request body into EntityType
			var e data.EntityType
			err = json.Unmarshal(body, &e)
			if err != nil {
				encodeError(api.RequestBodyParsingError, err, http.StatusBadRequest, w)
				return
			}

			// Persist parsed EntityType
			err = g.Service.Create(&e)
			if err != nil {
				encodeError(api.DatabasePersistingError, err, http.StatusInternalServerError, w)
				return
			}

			// Encode response
			encodeResponseBody(e, w)
		},
		ResponseHeaders: map[string]string{
			HeaderContentType: HeaderContentTypeValJSON,
		},
	}
}

// UpdateHandler returns a PUT endpoint handler for
// updaing the EntityType with the given id (path variable :id)
func (g *EntityTypeHandlerGroup) UpdateHandler() Handler {
	return Handler{
		Method: http.MethodPut,
		Path:   []string{strings.ToLower("EntityType"), ":id"},
		Logic: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Read request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				encodeError(api.RequestBodyReadingError, err, http.StatusBadRequest, w)
				return
			}

			// Parse request body into EntityType
			var e data.EntityType
			err = json.Unmarshal(body, &e)
			if err != nil {
				encodeError(api.RequestBodyParsingError, err, http.StatusBadRequest, w)
				return
			}

			// Persist parsed EntityType
			err = g.Service.Update(&e)
			if err != nil {
				encodeError(api.DatabasePersistingError, err, http.StatusInternalServerError, w)
				return
			}

			// Encode response
			encodeResponseBody(e, w)
		},
		ResponseHeaders: map[string]string{
			HeaderContentType: HeaderContentTypeValJSON,
		},
	}
}

// DeleteHandler returns a DELETE endpoint handler for
// deleting the EntityType at the given id (path variable :id)
func (g *EntityTypeHandlerGroup) DeleteHandler() Handler {
	return Handler{
		Method: http.MethodDelete,
		Path:   []string{strings.ToLower("EntityType"), ":id"},
		Logic: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Parse ID
			id, err := parsePathVarInteger(&ps, "id")
			if err != nil {
				encodeError(api.PathVariableParsingError, err, http.StatusBadRequest, w)
				return
			}

			// Delete by ID and retrieve existing
			e, err := g.Service.Delete(id)
			if err != nil {
				encodeError(api.DatabasePersistingError, err, http.StatusInternalServerError, w)
				return
			}

			// Encode response
			encodeResponseBody(e, w)

		},
		ResponseHeaders: map[string]string{
			HeaderContentType: HeaderContentTypeValJSON,
		},
	}
}

// GetAllHandler returns a GET endpoint handler for
// retrieving all EntityType
func (g *EntityTypeHandlerGroup) GetAllHandler() Handler {
	return Handler{
		Method: http.MethodGet,
		Path:   []string{strings.ToLower("EntityType")},
		Logic: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Retrieve all EntityType
			list, err := g.Service.GetAll()
			if err != nil {
				encodeError(api.DatabaseQueryingError, err, http.StatusInternalServerError, w)
				return
			}
			if list == nil {
				list = []data.EntityType{}
			}

			// Encode response
			encodeResponseBody(list, w)
		},
		ResponseHeaders: map[string]string{
			HeaderContentType: HeaderContentTypeValJSON,
		},
	}
}

// GetByIDHandler returns a GET endpoint handler for
// retrieving a single EntityType by the given id (path
// variable :id)
func (g *EntityTypeHandlerGroup) GetByIDHandler() Handler {
	return Handler{
		Method: http.MethodGet,
		Path:   []string{strings.ToLower("EntityType"), ":id"},
		Logic: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Parse ID from path
			id, err := parsePathVarInteger(&ps, "id")
			if err != nil {
				encodeError(api.PathVariableParsingError, err, http.StatusBadRequest, w)
				return
			}

			// Retrieve EntityType by ID
			e := data.EntityType{
				ID: id,
			}
			err = g.Service.GetByID(&e)
			if err != nil {
				encodeError(api.DatabaseQueryingError, err, http.StatusInternalServerError, w)
				return
			}

			// Encode response
			encodeResponseBody(e, w)
		},
		ResponseHeaders: map[string]string{
			HeaderContentType: HeaderContentTypeValJSON,
		},
	}
}
