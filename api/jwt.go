package api

import (
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/middle"
	"io"
	"net/http"
)

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	u, err := middle.UserFromJson(r.Body)
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
		return
	}

	err = dao.Authenticate(&u)
	if err != nil {
		handleHttpError(w, err, http.StatusUnauthorized)
		return
	}

	tokenString, err := middle.MakeToken(u)
	if err != nil {
		handleHttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}
