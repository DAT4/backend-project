package api

import (
	"fmt"
	"github.com/DAT4/backend-project/dao/mongo"
	"github.com/DAT4/backend-project/middle"
	"github.com/form3tech-oss/jwt-go"
	"io"
	"net/http"
	"time"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	u, err := middle.UserFromJson(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, err.Error())
		return
	}
	err = mongo.Authenticate(&u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, err.Error())
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": u.Id,
		"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(middle.AppKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"token_generation_failed"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}
