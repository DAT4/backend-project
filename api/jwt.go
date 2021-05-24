package api

import (
	"encoding/json"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/middle"
	"net/http"
)

type API struct {
	Db   dao.DBase
	Game *middle.Game
}

func (a *API) TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	u, err := middle.UserFromJson(r.Body)
	if err != nil {
		handleHttpError(w, "UserFromJson", err, http.StatusNotAcceptable)
		return
	}

	tokenPair, err := middle.AuthenticateUser(u, a.Db)
	if err != nil {
		handleHttpError(w, "Authenticate user", err, http.StatusUnauthorized)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tokenPair)
	if err != nil {
		handleHttpError(w, "EncodeTokenString", err, http.StatusInternalServerError)
	}
	return
}

func (a *API) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	token, err := middle.ExtractJWTToken(r, middle.REFRESH)
	if err != nil {
		return
	}
	u, err := middle.UserFromToken(token, a.Db)
	if err != nil {
		handleHttpError(w, "UserFromToken", err, http.StatusUnauthorized)
	}

	_, err = a.Db.UserFromName(string(u.Username))
	if err != nil {
		handleHttpError(w, "AuthenticateUser", err, http.StatusUnauthorized)
		return
	}

	tokenString, err := middle.RefreshToken(token, u)
	if err != nil {
		handleHttpError(w, "RefreshToken", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tokenString)
	if err != nil {
		handleHttpError(w, "EncodeTokenString", err, http.StatusInternalServerError)
	}
	return
}
