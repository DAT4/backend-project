package api

import (
	"github.com/DAT4/backend-project/models"
	"github.com/gorilla/mux"
)

func AddEndpoints(r *mux.Router) {
	endpoints := []models.Endpoint{
		{
			Path:    "/login",
			Handler: TokenHandler,
			Method:  "POST",
		},
		{
			Path:    "/register",
			Handler: CreateUser,
			Method:  "POST",
		},
		{
			Path:    "/join",
			Handler: JoinWebsocketConnection,
			Login:   true,
			Method:  "GET",
		},
	}
	for _, e := range endpoints {
		e.Add(r)
	}
}
