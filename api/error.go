package api

import (
	"fmt"
	"net/http"
)

func handleHttpError(w http.ResponseWriter, err error, status int) {
	fmt.Println(err)
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
