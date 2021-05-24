package middle

import (
	"errors"
	"fmt"
	"github.com/DAT4/backend-project/dto"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"net/http"
	"strings"
	"time"
)

const AppKey = "martin.mama.sh"

func MakeToken(u dto.User) (tokens TokenPair, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": u.Id,
		"exp":  time.Now().Add(time.Second * 15).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokens.AuthToken, err = token.SignedString([]byte(AppKey))
	if err != nil {
		return
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": u.Id,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokens.RefreshToken, err = token.SignedString([]byte(AppKey))
	return
}

type TokenPair struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

//https://medium.com/monstar-lab-bangladesh-engineering/jwt-auth-in-go-part-2-refresh-tokens-d334777ca8a0
func RefreshToken(refreshToken string, u dto.User) (tokens TokenPair, err error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppKey), nil
	})
	if err != nil {
		return
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokens, err = MakeToken(u)
	}
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

func extractClaims(tokenString string) (id string, err error) {
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

func ExtractJWTToken(req *http.Request, tokenType TokenType) (string, error) {
	var tokenString string
	switch tokenType {
	case AUTHENTICATION:
		tokenString = req.Header.Get("Authorization")
	case REFRESH:
		tokenString = req.Header.Get("RefreshToken")
	}
	if tokenString == "" {
		return "", fmt.Errorf("could not find token")
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
