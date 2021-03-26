package api

import (
	"github.com/DAT4/backend-project/dao/mongo"
	"github.com/DAT4/backend-project/middle"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	u , err := middle.UserFromJson(r.Body)
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
		return
	}
	err = middle.Validate(u)
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
		return
	}
	err = u.HashAndSalt()
	if err != nil {
		handleHttpError(w, err, http.StatusTeapot)
		return
	}
	err = mongo.Create(&u)
	if err != nil {
		handleHttpError(w, err, http.StatusTeapot)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}
