package api

import (
	"encoding/json"
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

	tmpUser, err := base.UserFromName(string(u.Username))
	if err != nil {
		handleHttpError(w, "AuthenticateUser", err, http.StatusUnauthorized)
		return
	}

	ok := u.Check(tmpUser.Password)
	if !ok {
		handleHttpError(w, "AuthenticateUser", err, http.StatusUnauthorized)
	}
	u.Id = tmpUser.Id
	u.PlayerID = tmpUser.PlayerID
	u.Password = tmpUser.Password

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

	_, err = base.UserFromName(string(u.Username))
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
