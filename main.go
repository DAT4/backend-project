package main

import (
	"fmt"
	"github.com/DAT4/backend-project/api"
	"github.com/DAT4/backend-project/middle"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	go middle.Game.Run()
	startREST()
}

func startREST() {
	r := mux.NewRouter()
	api.AddEndpoints(r)
	handler := cors.Default().Handler(r)
	fmt.Println("Running on port 8056")
	log.Fatal(http.ListenAndServe(":8056", handler))
}
