package api

import (
	"encoding/json"
	"fmt"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/middle"
	"net/http"
)

func tokenHandler(w http.ResponseWriter, r *http.Request, base dao.DBase) {
	w.Header().Add("Content-Type", "application/json")

	u, err := middle.UserFromJson(r.Body)
	if err != nil {
		handleHttpError(w, "UserFromJson", err, http.StatusNotAcceptable)
		return
	}

	err = base.Authenticate(&u)
	fmt.Println(u)
	if err != nil {
		handleHttpError(w, "AuthenticateUser", err, http.StatusUnauthorized)
		return
	}

	tokenString, err := middle.MakeToken(u)
	if err != nil {
		handleHttpError(w, "MakeToken", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tokenString)
	if err != nil {
		handleHttpError(w, "EncodeTokenString", err, http.StatusInternalServerError)
	}
	return
}

func refreshToken(w http.ResponseWriter, r *http.Request, base dao.DBase) {
	w.Header().Add("Content-Type", "application/json")
	token, err := middle.ExtractJWTToken(r, middle.REFRESH)
	if err != nil {
		return
	}
	u, err := middle.UserFromToken(token, base)
	if err != nil {
		handleHttpError(w, "UserFromToken", err, http.StatusUnauthorized)
	}

	err = base.Authenticate(&u)
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
