package middle

import (
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"net/http"
	"strings"
)

const AppKey = "martin.mama.sh"

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
