package models

import (
	"github.com/DAT4/backend-project/middle"
	"github.com/gorilla/mux"
	"net/http"
)

type Endpoint struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
	Login   bool
	Method  string
}

func (e Endpoint) Add(r *mux.Router) {
	if e.Login {
		r.Handle(e.Path, middle.AuthMiddleware(http.HandlerFunc(e.Handler))).Methods(e.Method)
	} else {
		r.HandleFunc(e.Path, e.Handler).Methods(e.Method)
	}
}

