package middle

import (
	"github.com/DAT4/backend-project/models/game"
	"github.com/DAT4/backend-project/models/user"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var Game *game.Game

func init() {
	Game = NewGame()
}

func NewGame() *game.Game {
	return &game.Game{
		Clients:    make(map[*game.Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *game.Client),
		Unregister: make(chan *game.Client),
	}
}

func ServeWs(u *user.User, g *game.Game, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &game.Client{User: u, Game: g, Conn: conn, Send: make(chan []byte, 256)}
	client.Game.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
