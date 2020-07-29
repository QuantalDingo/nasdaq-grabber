package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(username string, expirationTime time.Time) (string, error) {
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["expires_at"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("key")))

	return tokenString, err
}
