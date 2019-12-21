package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/friendsofgo/graphiql"
	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/handler"
	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/Dophin2009/nao/internal/data"
	"gitlab.com/Dophin2009/nao/internal/web"
)

// NewGraphQLHandler returns a POST endpoint handler for
// the GraphQL API.
func NewGraphQLHandler(ctx context.Context, schema *graphql.Schema, path []string) web.Handler {
	graphQLHandler := gqlhandler.New(&gqlhandler.Config{
		Schema:     schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: false,
	})
	return web.Handler{
		Method: http.MethodPost,
		Path:   path,
		Func: func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			graphQLHandler.ContextHandler(ctx, w, r)
		},
	}
}

// NewGraphiQLHandler returns a new GET endpoint handler
// for rendering a GraphiQL page for the given GraphQL API.
func NewGraphiQLHandler(path []string, graphqlPath string) (web.Handler, error) {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler(graphqlPath)
	if err != nil {
		return web.Handler{}, nil
	}
	return web.Handler{
		Method: http.MethodGet,
		Path:   path,
		Func: func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			graphiqlHandler.ServeHTTP(w, r)
		},
	}, nil
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
func LoginHandler(userService *data.UserService, jwtService *data.JWTService) web.Handler {
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
			err = userService.AuthenticateWithPassword(creds.Username, creds.Password)
			if err != nil {
				web.EncodeResponseErrorUnauthorized(web.ErrorAuthentication, err, w)
				return
			}

			user, err := userService.GetByUsername(creds.Username)

			expirationTime := time.Now().Add(5 * time.Minute)
			claims := &data.JWTClaims{
				UserID: user.ID,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			tokenString, err := jwtService.CreateTokenString(claims)
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
		ResponseHeaders: map[string]string{
			web.HeaderContentType: web.HeaderContentTypeValJSON,
		},
	}
}
