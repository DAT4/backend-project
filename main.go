package main

import (
	"flag"
	"fmt"
	"github.com/DAT4/backend-project/api"
	"github.com/DAT4/backend-project/dao/mongobase"
	"github.com/DAT4/backend-project/middle"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

// @title Backend API
// @version 0.5
// @description This is the api for the backend project
// @termsOfService https://backend.mama.sh/terms
// @contact.name Martin
// @contact.url https://mama.sh
// @contact.email mail@mama.sh
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host https://api.backend.mama.sh
func main() {
	var (
		addr string
		uri  string
	)

	flag.StringVar(&addr, "addr", ":8080", "The address and port to host on.")
	flag.StringVar(&uri, "db", "mongodb://localhost:27017", "The uri for the mongodb used.")
	flag.Parse()

	userDB := mongobase.NewUsers(&mongobase.Mongo{Uri: uri})

	server := api.API{
		Game: middle.NewGame(userDB),
	}

	go server.Game.Run()

	router := mux.NewRouter()

	users := router.PathPrefix("/user").Subrouter()
	{
		users.Path("").Methods(http.MethodGet).Handler(encrypt(server.FindUsers))
		users.Path("").Methods(http.MethodPost).HandlerFunc(server.InsertUser)
		users.Path("/{id}").Methods(http.MethodGet).Handler(encrypt(server.FindOneUser))
		users.Path("/{id}").Methods(http.MethodPut).Handler(encrypt(server.UpdateUser))
		users.Path("/{id}").Methods(http.MethodDelete).Handler(encrypt(server.DeleteUser))
	}

	router.Path("/login").Methods(http.MethodPost).HandlerFunc(server.TokenHandler)
	router.Path("/refresh").Methods(http.MethodPost).HandlerFunc(server.RefreshToken)
	router.Path("/join").Methods(http.MethodGet).HandlerFunc(server.JoinWebsocketConnection)

	handler := cors.Default().Handler(router)

	fmt.Printf("Running on port %v\n", addr)

	log.Fatal(http.ListenAndServe(addr, handler))
}

func encrypt(h http.HandlerFunc) http.Handler {
	return middle.AuthMiddleware(h)
}
