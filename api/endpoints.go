package api

import (
	"github.com/DAT4/backend-project/middle"
	"github.com/DAT4/backend-project/models"
	"github.com/gorilla/mux"
	"net/http"
)

func AddEndpoints(r *mux.Router) {
	endpoints := []models.Endpoint{
		{
			Path:    "/login",
			Handler: tokenHandler,
			Method:  "POST",
		},
		{
			Path:    "/register",
			Handler: createUser,
			Method:  "POST",
		},
		{
			Path:    "/join",
			Handler: joinWebsocketConnection,
			Method:  "GET",
		},
	}
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
