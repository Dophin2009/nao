package server

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/Dophin2009/nao/internal/data"
	"gitlab.com/Dophin2009/nao/internal/web"
	bolt "go.etcd.io/bbolt"
)

// MediaHandlerGroup is a basic handler group for Media.
type MediaHandlerGroup struct {
	Service           *data.MediaService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *MediaHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *MediaHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// EpisodeHandlerGroup is a basic handler group for Episode.
type EpisodeHandlerGroup struct {
	Service           *data.EpisodeService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *EpisodeHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *EpisodeHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *EpisodeHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *EpisodeHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *EpisodeHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *EpisodeHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *EpisodeHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// CharacterHandlerGroup is a basic handler group for Character.
type CharacterHandlerGroup struct {
	Service           *data.CharacterService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *CharacterHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *CharacterHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *CharacterHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *CharacterHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *CharacterHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *CharacterHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *CharacterHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// GenreHandlerGroup is a basic handler group for Character.
type GenreHandlerGroup struct {
	Service           *data.GenreService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *GenreHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *GenreHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *GenreHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *GenreHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *GenreHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *GenreHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *GenreHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// ProducerHandlerGroup is a basic handler group for Character.
type ProducerHandlerGroup struct {
	Service           *data.ProducerService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *ProducerHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *ProducerHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *ProducerHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *ProducerHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *ProducerHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *ProducerHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *ProducerHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// PersonHandlerGroup is a basic handler group for Character.
type PersonHandlerGroup struct {
	Service           *data.PersonService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *PersonHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *PersonHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *PersonHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *PersonHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *PersonHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *PersonHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *PersonHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// UserHandlerGroup is a basic handler group for User
type UserHandlerGroup struct {
	Service           *data.UserService
	JWTService        *data.JWTService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// LoginCredentials is passed in response body with
// the user's username and password to authenticate
type LoginCredentials struct {
	Username string
	Password string
}

// tokenCookieName is the name of the cookie the
// JWT token will be stored in in the login response
// and all subsequent requests
const tokenCookieName = "jwt_token"

// LoginHandler returns a POST endpoint handler to
// authenticate the user and return a JWT access
// token upon successful authentication
func (g *UserHandlerGroup) LoginHandler() web.Handler {
	return web.Handler{
		Method: http.MethodPost,
		Path:   []string{"login"},
		// On successful authentication, return JWT token
		Func: func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			// Read request body
			body, err := web.ReadRequestBody(r)
			if err != nil {
				web.EncodeResponseErrorBadRequest(web.ErrorRequestBodyReading, err, w)
				return
			}

			// Parse request body into login credentials
			var creds LoginCredentials
			err = json.Unmarshal(body, &creds)
			if err != nil {
				web.EncodeResponseErrorBadRequest(web.ErrorRequestBodyParsing, err, w)
				return
			}

			// Authenticate
			err = g.Service.AuthenticateWithPassword(creds.Username, creds.Password)
			if err != nil {
				web.EncodeResponseErrorUnauthorized(web.ErrorAuthentication, err, w)
				return
			}

			user := data.User{
				Username: creds.Username,
			}
			err = g.Service.GetByUsername(&user)

			expirationTime := time.Now().Add(5 * time.Minute)
			claims := &data.JWTClaims{
				UserID: user.ID,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			tokenString, err := g.JWTService.CreateTokenString(claims)
			if err != nil {
				web.EncodeResponseErrorInternalServer(web.ErrorInternalServer, err, w)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    tokenCookieName,
				Value:   tokenString,
				Expires: expirationTime,
			})
		},
	}
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *UserHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{
		g.LoginHandler(),
	}
}

func (g *UserHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadUsers:  true,
			WriteUsers: true,
		},
	)
}

func (g *UserHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadUsers:  true,
			WriteUsers: true,
		},
	)
}

func (g *UserHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadUsers:  true,
			WriteUsers: true,
		},
	)
}

func (g *UserHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadUsers: true,
		},
	)
}

func (g *UserHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadUsers: true,
		},
	)
}

func (g *UserHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// MediaRelationHandlerGroup is a handler group for MediaRelation.
type MediaRelationHandlerGroup struct {
	Service           *data.MediaRelationService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *MediaRelationHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *MediaRelationHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaRelationHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaRelationHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaRelationHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaRelationHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaRelationHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// MediaCharacterHandlerGroup is a handler group for MediaRelation.
type MediaCharacterHandlerGroup struct {
	Service           *data.MediaCharacterService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *MediaCharacterHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *MediaCharacterHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaCharacterHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaCharacterHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaCharacterHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaCharacterHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaCharacterHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// MediaGenreHandlerGroup is a handler group for MediaRelation.
type MediaGenreHandlerGroup struct {
	Service           *data.MediaGenreService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *MediaGenreHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *MediaGenreHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaGenreHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaGenreHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaGenreHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaGenreHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaGenreHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// MediaProducerHandlerGroup is a handler group for MediaRelation.
type MediaProducerHandlerGroup struct {
	Service           *data.MediaProducerService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *MediaProducerHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *MediaProducerHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaProducerHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaProducerHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			WriteMedia: true,
		},
	)
}

func (g *MediaProducerHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaProducerHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
		},
	)
}

func (g *MediaProducerHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// UserMediaHandlerGroup is a handler group for MediaRelation.
type UserMediaHandlerGroup struct {
	Service           *data.UserMediaService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *UserMediaHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *UserMediaHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			ReadUsers:  true,
			WriteMedia: true,
			WriteUsers: true,
		},
	)
}

func (g *UserMediaHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			ReadUsers:  true,
			WriteMedia: true,
			WriteUsers: true,
		},
	)
}

func (g *UserMediaHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			ReadUsers:  true,
			WriteMedia: true,
			WriteUsers: true,
		},
	)
}

func (g *UserMediaHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
			ReadUsers: true,
		},
	)
}

func (g *UserMediaHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
			ReadUsers: true,
		},
	)
}

func (g *UserMediaHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// UserMediaListHandlerGroup is a handler group for MediaRelation.
type UserMediaListHandlerGroup struct {
	Service           *data.UserMediaListService
	JWTAuthenticator  *JWTAuthenticator
	PermAuthenticator *UserPermissionAuthenticator
}

// ExtraHandlers returns extra handlers besides the CRUD
// handlers for the handler group.
func (g *UserMediaListHandlerGroup) ExtraHandlers() []web.Handler {
	return []web.Handler{}
}

func (g *UserMediaListHandlerGroup) createAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			ReadUsers:  true,
			WriteMedia: true,
			WriteUsers: true,
		},
	)
}

func (g *UserMediaListHandlerGroup) updateAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			ReadUsers:  true,
			WriteMedia: true,
			WriteUsers: true,
		},
	)
}

func (g *UserMediaListHandlerGroup) deleteAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia:  true,
			ReadUsers:  true,
			WriteMedia: true,
			WriteUsers: true,
		},
	)
}

func (g *UserMediaListHandlerGroup) getAllAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
			ReadUsers: true,
		},
	)
}

func (g *UserMediaListHandlerGroup) getByIDAuthenticator(r *http.Request, ps httprouter.Params) error {
	return g.authenticateJWTAndPermissions(
		r, ps, &data.Permission{
			ReadMedia: true,
			ReadUsers: true,
		},
	)
}

func (g *UserMediaListHandlerGroup) authenticateJWTAndPermissions(r *http.Request, ps httprouter.Params, perm *data.Permission) error {
	return authenticateJWTAndPermissions(r, ps, g.JWTAuthenticator, g.PermAuthenticator, perm)
}

// NewEntityHandlerGroups returns a list of all
// entity-related handler groups with
func NewEntityHandlerGroups(db *bolt.DB) []web.HandlerGroup {
	return []web.HandlerGroup{
		&MediaHandlerGroup{
			Service: &data.MediaService{DB: db},
		},
		&EpisodeHandlerGroup{
			Service: &data.EpisodeService{DB: db},
		},
		&CharacterHandlerGroup{
			Service: &data.CharacterService{DB: db},
		},
		&GenreHandlerGroup{
			Service: &data.GenreService{DB: db},
		},
		&ProducerHandlerGroup{
			Service: &data.ProducerService{DB: db},
		},
		&PersonHandlerGroup{
			Service: &data.PersonService{DB: db},
		},
		&UserHandlerGroup{
			Service: &data.UserService{DB: db},
		},
		&MediaRelationHandlerGroup{
			Service: &data.MediaRelationService{DB: db},
		},
		&MediaCharacterHandlerGroup{
			Service: &data.MediaCharacterService{DB: db},
		},
		&MediaGenreHandlerGroup{
			Service: &data.MediaGenreService{DB: db},
		},
		&MediaProducerHandlerGroup{
			Service: &data.MediaProducerService{DB: db},
		},
		&UserMediaHandlerGroup{
			Service: &data.UserMediaService{DB: db},
		},
		&UserMediaListHandlerGroup{
			Service: &data.UserMediaListService{DB: db},
		},
	}
}

func authenticateJWTAndPermissions(
	r *http.Request, ps httprouter.Params,
	jwtAuthenticator *JWTAuthenticator,
	permAuthenticator *UserPermissionAuthenticator,
	req *data.Permission,
) error {
	claims, err := jwtAuthenticator.Authenticate(r, ps)
	if err != nil {
		return err
	}

	_, err = permAuthenticator.Authenticate(claims.UserID, req)
	return err
}
