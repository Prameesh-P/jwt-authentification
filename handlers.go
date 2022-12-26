package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var JwtKey = []byte("secret_key")

// Users this just a sample to store username and password.You can use the help of database.

var Users = map[string]string{
	"Pramee":   "123",
	"Prameesh": "12345",
}

// Credentials for check api request data..
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct {
	Username           string `json:"username"`
	jwt.StandardClaims        //Inside this have lot of preset value.Now I am using Expire At keyword.
}

// Login this is a login handler
func Login(w http.ResponseWriter, r *http.Request) {
	// here we are create the object of the Credentials struct.
	var credentials Credentials
	//and decode that object
	err := json.NewDecoder(r.Body).Decode(&credentials)
	// if any error handle.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//Nothing on the credential return to badRequest
		return
	}
	//if anything on that credential we want check our expected password and credential password is same.
	expectedPassword, ok := Users[credentials.Username] //if you get actual password ok is true.
	//else ok is false
	if !ok || expectedPassword != credentials.Password { // not ok(!ok) means our expected password is wrong

		w.WriteHeader(http.StatusBadRequest) // and return to badRequest
		return
	}
	//everything is ok we are creating an expiration time
	expirationTime := time.Now().Add(time.Minute * 5)
	//and we are updating Claims variable
	claims := Claims{
		Username:       credentials.Username,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims) //and finally we are passing signing method.
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} //and setting our cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

// Home this is a home handler
func Home(w http.ResponseWriter, r *http.Request) {
	//we are passing our cookie value "token".
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			//if no cookie that is Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//if any other error that is a bad request
		w.WriteHeader(http.StatusBadRequest)
	}
	//storing our cookie value to tokenStr
	tokenStr := cookie.Value
	claims := &Claims{} //and storing Claims on claims variable
	//we are passing with jwt.ParseWithClaims
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid { //if token not valid status unauthorized
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//else we are welcome the user
	w.Write([]byte(fmt.Sprintf("Hello %s", claims.Username)))
}

// Refresh and we want to refresh our token.
func Refresh(w http.ResponseWriter, r *http.Request) {
	//this is similar to Home handler
	//we are passing our cookie value "token".
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			//if no cookie that is Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//if any other error that is a bad request
		w.WriteHeader(http.StatusBadRequest)
	}
	//storing our cookie value to tokenStr
	tokenStr := cookie.Value
	claims := &Claims{} //and storing Claims on claims variable
	//we are passing with jwt.ParseWithClaims
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid { //if token not valid status unauthorized
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//30 second left on our old token and that time we will create new token
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		//we are checking the token expiration time is greater than 30 second
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//this all method similar to the Login handler
	expirationTime := time.Now().Add(time.Minute * 5)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "new_token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
