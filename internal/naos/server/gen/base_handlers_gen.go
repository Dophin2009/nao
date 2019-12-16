package server

import (
	"net/http"
	"strings"

	"github.com/cheekybits/genny/generic"
	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/Dophin2009/nao/internal/data"
	"gitlab.com/Dophin2009/nao/internal/web"
)

//go:generate genny -in=base_handlers_gen.go -out=base_handlers.gen.go gen "EntityType=Media,Episode,Character,Genre,Producer,Person,User,MediaRelation,MediaCharacter,MediaGenre,MediaProducer,UserMedia,UserMediaList"

// EntityType is a generic placeholder for all entity types;
// it is assumed that EntityType structs have an ID and Version,
// both of which should be of int type.
type EntityType generic.Type

// Handlers returns all the handlers for the handler group
func (g *EntityTypeHandlerGroup) Handlers() []web.Handler {
	handlers := []web.Handler{
		g.CreateHandler(),
		g.UpdateHandler(),
		g.DeleteHandler(),
		g.GetAllHandler(),
		g.GetByIDHandler(),
	}
	handlers = append(handlers, g.ExtraHandlers()...)
	return handlers
}

// CreateHandler returns an POST endpoint handler
// for creating new EntityType
func (g *EntityTypeHandlerGroup) CreateHandler() web.Handler {
	return web.Handler{
		Method: http.MethodPost,
		Path:   []string{strings.ToLower("EntityType")},
		Func: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			err := g.createAuthenticator(r, ps)
			if err != nil {
				web.EncodeResponseErrorUnauthorized(web.ErrorAuthentication, err, w)
				return
			}

			// Read request body
			body, err := web.ReadRequestBody(r)
			if err != nil {
				web.EncodeResponseErrorBadRequest(web.ErrorRequestBodyReading, err, w)
				return
			}

			// Parse request body into EntityType
			var e data.EntityType
			err = json.Unmarshal(body, &e)
			if err != nil {
				web.EncodeResponseErrorBadRequest(web.ErrorRequestBodyParsing, err, w)
				return
			}

			// Persist parsed EntityType
			err = g.Service.Create(&e)
			if err != nil {
				web.EncodeResponseErrorInternalServer(web.ErrorInternalServer, err, w)
				return
			}

			// Encode response
			web.EncodeResponseBody(e, w)
		},
		ResponseHeaders: map[string]string{
			web.HeaderContentType: web.HeaderContentTypeValJSON,
		},
	}
}

// UpdateHandler returns a PUT endpoint handler for
// updaing the EntityType with the given id (path variable :id)
func (g *EntityTypeHandlerGroup) UpdateHandler() web.Handler {
	return web.Handler{
		Method: http.MethodPut,
		Path:   []string{strings.ToLower("EntityType"), ":id"},
		Func: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			err := g.updateAuthenticator(r, ps)
			if err != nil {
				web.EncodeResponseErrorUnauthorized(web.ErrorAuthentication, err, w)
				return
			}

			// Read request body
			body, err := web.ReadRequestBody(r)
			if err != nil {
				web.EncodeResponseErrorBadRequest(web.ErrorRequestBodyReading, err, w)
				return
			}

			// Parse request body into EntityType
			var e data.EntityType
			err = json.Unmarshal(body, &e)
			if err != nil {
				web.EncodeResponseErrorBadRequest(web.ErrorRequestBodyParsing, err, w)
				return
			}

			// Persist parsed EntityType
			err = g.Service.Update(&e)
			if err != nil {
				web.EncodeResponseErrorInternalServer(web.ErrorInternalServer, err, w)
				return
			}

			// Encode response
			web.EncodeResponseBody(e, w)
		},
		ResponseHeaders: map[string]string{
			web.HeaderContentType: web.HeaderContentTypeValJSON,
		},
	}
}

// DeleteHandler returns a DELETE endpoint handler for
// deleting the EntityType at the given id (path variable :id)
func (g *EntityTypeHandlerGroup) DeleteHandler() web.Handler {
	return web.Handler{
		Method: http.MethodDelete,
		Path:   []string{strings.ToLower("EntityType"), ":id"},
		Func: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			err := g.deleteAuthenticator(r, ps)
			if err != nil {
				web.EncodeResponseErrorUnauthorized(web.ErrorAuthentication, err, w)
				return
			}

			// Parse ID
			id, err := web.ParsePathVarInt("id", &ps)
			if err != nil {
				web.EncodeResponseErrorBadRequest(web.ErrorPathVariableParsing, err, w)
				return
			}

			// Delete by ID and retrieve existing
			e, err := g.Service.Delete(id)
			if err != nil {
				web.EncodeResponseErrorInternalServer(web.ErrorInternalServer, err, w)
				return
			}

			// Encode response
			web.EncodeResponseBody(e, w)
		},
		ResponseHeaders: map[string]string{
			web.HeaderContentType: web.HeaderContentTypeValJSON,
		},
	}
}

// GetAllHandler returns a GET endpoint handler for
// retrieving all EntityType
func (g *EntityTypeHandlerGroup) GetAllHandler() web.Handler {
	return web.Handler{
		Method: http.MethodGet,
		Path:   []string{strings.ToLower("EntityType")},
		Func: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			err := g.getAllAuthenticator(r, ps)
			if err != nil {
				web.EncodeResponseErrorUnauthorized(web.ErrorAuthentication, err, w)
				return
			}

			// Retrieve all EntityType
			list, err := g.Service.GetAll()
			if err != nil {
				web.EncodeResponseErrorInternalServer(web.ErrorInternalServer, err, w)
				return
			}
			if list == nil {
				list = []data.EntityType{}
			}

			// Encode response
			web.EncodeResponseBody(list, w)
		},
		ResponseHeaders: map[string]string{
			web.HeaderContentType: web.HeaderContentTypeValJSON,
		},
	}
}

// GetByIDHandler returns a GET endpoint handler for
// retrieving a single EntityType by the given id (path
// variable :id)
func (g *EntityTypeHandlerGroup) GetByIDHandler() web.Handler {
	return web.Handler{
		Method: http.MethodGet,
		Path:   []string{strings.ToLower("EntityType"), ":id"},
		Func: func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			err := g.getByIDAuthenticator(r, ps)
			if err != nil {
				web.EncodeResponseErrorUnauthorized(web.ErrorAuthentication, err, w)
				return
			}

			// Parse ID from path
			id, err := web.ParsePathVarInt("id", &ps)
			if err != nil {
				web.EncodeResponseErrorBadRequest(web.ErrorPathVariableParsing, err, w)
				return
			}

			// Retrieve EntityType by ID
			e := data.EntityType{
				ID: id,
			}
			err = g.Service.GetByID(&e)
			if err != nil {
				web.EncodeResponseErrorInternalServer(web.ErrorInternalServer, err, w)
				return
			}

			// Encode response
			web.EncodeResponseBody(e, w)
		},
		ResponseHeaders: map[string]string{
			web.HeaderContentType: web.HeaderContentTypeValJSON,
		},
	}
}
