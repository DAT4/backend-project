package main

import (
	"backend-projekt/models"
	"errors"
	"fmt"
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
			games[newGame.Name] = newGame //DATA RACE DETECTED
			go newGame.Run()
		}
	}
}

func main() {
	go runHub()
	r := mux.NewRouter()
	r.HandleFunc("/login", models.TokenHandler).Methods("POST")
	r.HandleFunc("/register", createUser).Methods("POST")
	r.Handle("/create", models.AuthMiddleware(http.HandlerFunc(CreateWebsocketConnection))).Methods("GET")
	r.Handle("/join", models.AuthMiddleware(http.HandlerFunc(JoinWebsocketConnection))).Methods("GET")
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
