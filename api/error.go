package api

import (
	"fmt"
	"net/http"
)

func handleHttpError(w http.ResponseWriter, context string, err error, status int) {
	fmt.Println(err)
	w.WriteHeader(status)
	w.Write([]byte(context + ": " + err.Error()))
}
