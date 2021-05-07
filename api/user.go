package api

import (
	"encoding/json"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/middle"
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request, base dao.DBase) {
	w.Header().Add("Content-Type", "application/json")

	u, err := middle.UserFromJson(r.Body)
	if err != nil {
		handleHttpError(w, "UserFromJson", err, http.StatusNotAcceptable)
		return
	}

	u, err = middle.CreateUser(u, base)
	if err != nil {
		handleHttpError(w, "ValidateUser", err, http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		handleHttpError(w, "Serializing json", err, http.StatusInternalServerError)
		return
	}
}
