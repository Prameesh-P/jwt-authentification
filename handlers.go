package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var JwtKey = []byte("secret_key")
var Users = map[string]string{
	"Pramee":   "123",
	"Prameesh": "12345",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Home(w http.ResponseWriter, r *http.Request) {

}
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exepectedPassword, ok := Users[credentials.Username]
	if !ok || exepectedPassword != credentials.Password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := Claims{
		Username: credentials.Username,
		StandardClaims:jwt.StandardClaims{ExpiresAt: expirationTime.Unix()
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodES256,claims)
	tokenString,err:=token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w,&http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})

}
func Refresh(w http.ResponseWriter, r *http.Request) {

}
