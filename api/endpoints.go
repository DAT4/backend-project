package api

import (
	"github.com/DAT4/backend-project/models"
	"github.com/gorilla/mux"
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
			Login:   true,
			Method:  "GET",
		},
	}
	for _, e := range endpoints {
		e.Add(r)
	}
}
