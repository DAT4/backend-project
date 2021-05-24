package main

import (
	"flag"
	"fmt"
	"github.com/DAT4/backend-project/api"
	"github.com/DAT4/backend-project/dao"
	"github.com/DAT4/backend-project/middle"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	var (
		addr string
		uri  string
	)

	flag.StringVar(&addr, "addr", ":8080", "The address and port to host on.")
	flag.StringVar(&uri, "db", "mongodb://localhost:27017", "The uri for the mongodb used.")
	flag.Parse()

	db, err := dao.NewMongoDB(uri)
	if err != nil {
		return
	}

	server := api.API{
		Game: middle.NewGame(db),
	}

	go server.Game.Run()

	router := mux.NewRouter()

	router.Path("/register").Methods(http.MethodPost).HandlerFunc(server.CreateUser)
	router.Path("/login").Methods(http.MethodPost).HandlerFunc(server.TokenHandler)
	router.Path("/refresh").Methods(http.MethodPost).HandlerFunc(server.RefreshToken)
	router.Path("/join").Methods(http.MethodGet).HandlerFunc(server.JoinWebsocketConnection)

	handler := cors.Default().Handler(router)

	fmt.Printf("Running on port %v\n", addr)

	log.Fatal(http.ListenAndServe(addr, handler))
}
