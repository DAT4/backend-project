package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Data  string `json:"data" example:"some message"`
	Error string `json:"error" example:"some error message"`
}

func handleError(tag string, err error, w http.ResponseWriter, status int) {
	log.Println("HTTP request had an error:", err.Error())
	resp := Response{Error: fmt.Sprintf("%s: %s", tag, err)}
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("An error occurred while returning error:", err)
	}
}
