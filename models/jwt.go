package models

import (
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"io"
	"net/http"
	"strings"
	"time"
)

const AppKey = "martin.mama.sh"

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	user, err := UserFromJson(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, err.Error())
		return
	}
	user, err = user.Authenticate()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, err.Error())
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Id,
		"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(AppKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"token_generation_failed"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}

func AuthMiddleware(next http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(AppKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}
func ExtractClaims(tokenString string) (id string, err error) {
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(AppKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}
	id, ok := claims["user"].(string)
	if !ok {
		return id, errors.New("no user in map")
	}
	return
}

func ExtractJWTToken(req *http.Request) (string, error) {
	tokenString := req.Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("Could not find token")
	}
	tokenString, err := stripTokenPrefix(tokenString)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func stripTokenPrefix(tok string) (string, error) {
	tokenParts := strings.Split(tok, " ")
	if len(tokenParts) < 2 {
		return tokenParts[0], nil
	}
	return tokenParts[1], nil
}
