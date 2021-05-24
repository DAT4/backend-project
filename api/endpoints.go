package api

import (
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/middle"
	"github.com/DAT4/backend-project/models"
	"github.com/gorilla/mux"
	"net/http"
)

func byLazy(fn func(w http.ResponseWriter, r *http.Request, base dao.DBase), db dao.DBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func AddEndpoints(r *mux.Router, db dao.DBase) {
	endpoints := []models.Endpoint{}
	add(endpoints, r)
}

func add(es []models.Endpoint, r *mux.Router) {
	for _, e := range es {
		if e.Login {
			r.Handle(e.Path, middle.AuthMiddleware(http.HandlerFunc(e.Handler))).Methods(e.Method)
		} else {
			r.HandleFunc(e.Path, e.Handler).Methods(e.Method)
		}
	}
}
