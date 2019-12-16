package server

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/Dophin2009/nao/internal/data"
	"gitlab.com/Dophin2009/nao/internal/web"
)

// JWTAuthenticator is an authenticator for JWT tokens.
type JWTAuthenticator struct {
	Service *data.JWTService
}

// Authenticate checks for the presence and validity of a
// token cookie in the request.
func (au *JWTAuthenticator) Authenticate(r *http.Request, _ httprouter.Params) (claims data.JWTClaims, err error) {
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

// UserPermissionAuthenticator is an authenticator that
// checks whether the user has sufficient permissions.
type UserPermissionAuthenticator struct {
	Service *data.UserService
}

// Authenticate checks if the user with the given ID
// has permissions that meet the requirements.
func (au *UserPermissionAuthenticator) Authenticate(userID int, req *data.Permission) (user data.User, err error) {
	user.ID = userID
	err = au.Service.GetByID(&user)
	if err != nil {
		return
	}

	if !au.RequirementsMet(&user.Permissions, req) {
		err = errors.New("insufficient permissions")
		return
	}
	return
}

// RequirementsMet checks if the given permissions satisfy
// the requiremed permissions.
func (au *UserPermissionAuthenticator) RequirementsMet(perm *data.Permission, req *data.Permission) bool {
	return !((req.ReadMedia && !perm.ReadMedia) || (req.WriteMedia && !perm.WriteMedia))
}
