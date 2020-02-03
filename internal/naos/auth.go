package naos

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Dophin2009/nao/internal/data"
	"gitlab.com/Dophin2009/nao/internal/web"
)

// JWTAuthenticator is an authenticator for JWT tokens.
type JWTAuthenticator struct {
	Service *data.JWTService
}

// Authenticate checks for the presence and validity of a token cookie in the
// request.
func (au *JWTAuthenticator) Authenticate(
	r *http.Request, _ httprouter.Params) (claims data.JWTClaims, err error) {
	c, err := r.Cookie(tokenCookieName)
	if err != nil {
		return
	}

	tokenString := c.Value
	claims, err = au.Service.CheckTokenString(tokenString)
	if err != nil {
		err = &web.AuthenticationError{
			Debug: err.Error(),
		}
		return
	}

	return
}
