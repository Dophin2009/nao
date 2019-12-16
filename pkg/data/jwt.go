package data

import (
	"errors"
	"math/rand"

	"github.com/dgrijalva/jwt-go"
	bolt "go.etcd.io/bbolt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// JWTService performs operations on the peristence
// layer
type JWTService struct {
	DB        *bolt.DB
	secretKey *string
}

// JWTBucket is the name of the bucket to store JWT
// secret keys
const JWTBucket = "JWT"

type JWTClaims struct {
	UserID int
	jwt.StandardClaims
}

// CreateTokenString returns an encoded token string
// from the given claims
func (ser *JWTService) CreateTokenString(claims *JWTClaims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := ser.secret()
	if err != nil {
		return
	}
	tokenString, err = token.SignedString(secret)
	return
}

// CheckTokenString checks if the token string is valid.
func (ser *JWTService) CheckTokenString(tokenString string) (claims JWTClaims, err error) {
	claims = JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return ser.secret()
	})
	if err != nil {
		return
	}

	if !token.Valid {
		err = errors.New("token invalid")
	}
	return
}

// GenerateNewSecret returns a new randomly generated
// string of the given length.
func (ser *JWTService) generateNewSecret(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// secret returns the cached JWT secret key,
// or retrieves the secret from the database if it
// has not been cached.
func (ser *JWTService) secret() (secret string, err error) {
	// Check if secret has already been cached in service
	if ser.secretKey != nil {
		secret = *ser.secretKey
		return
	}
	// If not cached, check database for secret
	err = ser.DB.Update(func(tx *bolt.Tx) error {
		// Retrieve secret from database
		b := tx.Bucket([]byte(JWTBucket))
		v := b.Get([]byte("secret"))
		// If secret is not in database, generate a new one
		// and store
		if v == nil {
			secret = ser.generateNewSecret(15)
			ser.secretKey = &secret
			err = b.Put([]byte("secret"), []byte(secret))
			if err != nil {
				return err
			}
		} else {
			secret = string(v)
			ser.secretKey = &secret
		}
		return nil
	})
	return
}
