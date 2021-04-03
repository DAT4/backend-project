package api

import (
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/middle"
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	u, err := middle.UserFromJson(r.Body)
	if err != nil {
		handleHttpError(w, "UserFromJson", err, http.StatusNotAcceptable)
		return
	}
	err = middle.Validate(u)
	if err != nil {
		handleHttpError(w, "ValidateUser", err, http.StatusNotAcceptable)
		return
	}
	err = u.HashAndSalt()
	if err != nil {
		handleHttpError(w, "HashAndSalt", err, http.StatusTeapot)
		return
	}
	err = dao.Create(&u)
	if err != nil {
		handleHttpError(w, "CreateUser", err, http.StatusTeapot)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}
