package domain

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

const issuer string = "mainflux"

var secretKey string = "mainflux-api-key"

type tokenClaims struct {
	Payload
	jwt.StandardClaims
}

// SetSecretKey sets the secret key that will be used for decoding and encoding
// generated tokens. If not invoked, a default key will be used instead.
func SetSecretKey(key string) {
	secretKey = key
}

// CreateKey creates new JSON Web Token containing provided subject and
// payload.
func CreateKey(subject string, payload *Payload) (string, error) {
	claims := tokenClaims{
		*payload,
		jwt.StandardClaims{
			Issuer:  issuer,
			Subject: subject,
		},
	}

	key := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := key.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return raw, nil
}

// ScopesOf extracts the token payload given its string representation.
func ScopesOf(key string) (Payload, error) {
	var payload Payload

	token, err := jwt.ParseWithClaims(
		key,
		&tokenClaims{},
		func(val *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return payload, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims.Payload, nil
	}

	return payload, errors.New("failed to retrieve claims")
}
