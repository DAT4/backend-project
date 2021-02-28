package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	fmt.Println("User inside")
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
	err = user.hashAndSalt()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte(err.Error()))
		return
	}
	var tmpUser User
	q1 := FindOneQuery{
		Model:      &tmpUser,
		Filter:     bson.M{"username": user.Username},
		Options:    options.FindOne(),
		Collection: "users",
	}
	err = q1.find()
	if err == nil {
		fmt.Println("Finding: ",err)
		fmt.Println(tmpUser)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("A user already exists with this name"))
		return
	}
	q2 := AddOneQuery{
		Model:      &user,
		Filter:     nil,
		Collection: "users",
	}

	err = q2.add()
	if err != nil {
		fmt.Println("Adding: ",err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}
