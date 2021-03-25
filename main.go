package main

import (
	"errors"
	"fmt"
	"github.com/DAT4/backend-project/models"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

var games = make(map[models.Username]*models.Game) //DATA RACE DETECTED
var gamesChannel = make(chan *models.Game)

func runHub() {
	for {
		select {
		case newGame := <-gamesChannel:
			fmt.Println("Game started by", newGame.Name)
			games[newGame.Name] = newGame //DATA RACE DETECTED
			go newGame.Run()
		}
	}
}

type endpoint struct {
	path string
	handler func(w http.ResponseWriter, r *http.Request)
	secure bool
	method string
}

func main() {
	go runHub()
	r := mux.NewRouter()
	endpoints := []endpoint{
		{"/login",models.TokenHandler, false,"POST" },
		{"/register",createUser, false,"POST" },
		{"/create",CreateWebsocketConnection, true,"GET" },
		{"/join",JoinWebsocketConnection, true,"GET" },
	}

	for _, e := range endpoints{
		if e.secure {
			r.Handle(e.path, models.AuthMiddleware(http.HandlerFunc(e.handler))).Methods(e.method)
		} else {
			r.HandleFunc(e.path, e.handler).Methods(e.method)
		}

	}

	handler := cors.Default().Handler(r)
	fmt.Println("Running on port 8056")
	log.Fatal(http.ListenAndServe(":8056", handler))
}

func CreateWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	user, err := models.UserFromToken(r)
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
	}
	game := models.NewGame(user.Username)
	gamesChannel <- game
	user.ServeWs(game, w, r)
}

func JoinWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Peer")
	user, err := models.UserFromToken(r)
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
	}
	peer := r.URL.Query().Get("peer")
	if len(peer) > 0 {
		game, ok := games[models.Username(peer)]
		if !ok {
			err = errors.New("peer available but no Game is found")
			handleHttpError(w, err, http.StatusInternalServerError)
			return
		}
		user.ServeWs(game, w, r) //DATA RACE DETECTED
	} else {
		err = errors.New("peer not available")
		handleHttpError(w, err, http.StatusUnavailableForLegalReasons)
	}
}

func handleHttpError(w http.ResponseWriter, err error, status int) {
	fmt.Println(err)
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	user, err := models.UserFromJson(r.Body)
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
		return
	}
	err = user.Validate()
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
		return
	}
	err = user.HashAndSalt()
	if err != nil {
		handleHttpError(w, err, http.StatusTeapot)
		return
	}
	err = user.Create()
	if err != nil {
		handleHttpError(w, err, http.StatusTeapot)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}
