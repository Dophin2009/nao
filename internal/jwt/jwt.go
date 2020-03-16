package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

const keyEnvKey = "JWT_KEY"

// Authenticator authenticates JSON web tokens.
type Authenticator struct {
	key string
}

// Claims is a custom JWT claims type with username and expiration information.
type Claims struct {
	Username string
	jwt.StandardClaims
}

// Verify checks the given HTTP request for a valid JWT.
func (au *Authenticator) Verify(tokenstr string) error {
	claims := Claims{}
	tkn, err := jwt.ParseWithClaims(tokenstr, &claims,
		func(_ *jwt.Token) (interface{}, error) {
			return au.key, nil
		})
	if err != nil {
		return fmt.Errorf("failed to parse token string: %w", err)
	}

	if !tkn.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

// NewToken returns a new JWT token.
func (au *Authenticator) NewToken(username string, minDuration time.Duration) (string, error) {
	expiration := time.Now().Add(minDuration * time.Minute)
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tknstr, err := token.SignedString(au.key)
	if err != nil {
		return "", fmt.Errorf("failed to create signed string: %w", err)
	}
	return tknstr, nil
}

// ReadKeyFromEnv reads the JWT secret key from a .env file at the given path
// and returns it.
func ReadKeyFromEnv(filepath string) (string, error) {
	err := godotenv.Load(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to load env file %q: %w", filepath, err)
	}

	key := os.Getenv(keyEnvKey)
	return key, nil
}

// WriteKeyToEnv writes the given JWT secret key to a .env file at the given
// path.
func WriteKeyToEnv(key string, filepath string) error {
	env, err := godotenv.Unmarshal(fmt.Sprintf("%s=%s", keyEnvKey, key))
	if err != nil {
		return fmt.Errorf("failed to unmarshal env: %w", err)
	}

	err = godotenv.Write(env, filepath)
	if err != nil {
		return fmt.Errorf("failed to write env to file: %w", err)
	}
	return nil
}
