package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"nasdaq-grabber/models"

	"github.com/dgrijalva/jwt-go"
)

func SignIn(rw http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = user.Create()
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	MakeToken(rw, user.Username)
}

func LogIn(rw http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.IsExist() && user.Validate() {
		MakeToken(rw, user.Username)
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

func Refresh(rw http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("key")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	MakeToken(rw, fmt.Sprintf("%v", claims["username"]))
}

func MakeToken(rw http.ResponseWriter, username string) {
	expirationTime := time.Now().Add(15 * time.Minute)

	token, err := models.CreateToken(username, expirationTime)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	setTokenToCookie(rw, token, expirationTime)
}

func setTokenToCookie(rw http.ResponseWriter, token string, expirationTime time.Time) {
	http.SetCookie(rw, &http.Cookie{
		Name:    "access_token",
		Value:   token,
		Expires: expirationTime,
	})
}
