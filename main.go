package main

import (
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
	db := &dao.MongoDB{}
	go middle.G.Run(db)
	startREST(db)
}

func startREST(db dao.DBase) {
	r := mux.NewRouter()
	api.AddEndpoints(r, db)
	handler := cors.Default().Handler(r)
	fmt.Println("Running on port 1001")
	log.Fatal(http.ListenAndServe(":1001", handler))
}
