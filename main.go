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
		db   string
	)
	flag.StringVar(&addr, "addr", ":8080", "The address and port to host on.")
	flag.StringVar(&db, "db", "mongodb://localhost:27017", "The uri for the mongodb used.")
	flag.Parse()
	startREST(db, addr)
}

func startREST(dbURI, addr string) {
	fmt.Println(dbURI)
	db, err := dao.NewMongoDB(dbURI)
	if err != nil {
		log.Fatalf("Could not create the database: %v", err)
	}

	go middle.G.Run(&db)
	r := mux.NewRouter()
	api.AddEndpoints(r, &db)
	handler := cors.Default().Handler(r)
	fmt.Printf("Running on port %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
