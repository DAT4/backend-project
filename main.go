package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", TokenHandler).Methods("POST")
	r.HandleFunc("/register", createUser).Methods("POST")
	r.Handle("/game", AuthMiddleware(http.HandlerFunc(Game))).Methods("POST")
	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":8055", handler))
}

func Game(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Hope you like tea"))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user User
	err := user.fromJson(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(err.Error()))
		return
	}
	err = user.validate()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(err.Error()))
		return
	}
}
