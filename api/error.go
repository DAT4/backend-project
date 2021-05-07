package api

import (
	"net/http"
)

func handleHttpError(w http.ResponseWriter, context string, err error, status int) {
	var errmsg string
	w.WriteHeader(status)
	if err != nil {
		errmsg = err.Error()
	}
	w.Write([]byte(context + ": " + errmsg))
}
