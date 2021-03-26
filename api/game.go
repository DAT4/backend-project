package api

import (
	"fmt"
	"github.com/DAT4/backend-project/middle"
	"github.com/DAT4/backend-project/models"
	"github.com/DAT4/backend-project/models/game"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func joinWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Peer")
	u, err := middle.UserFromToken(r)
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
	}
	serveWs(&u, middle.Game, w, r)
}

func serveWs(u *models.User, g *game.Game, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	game.NewClient(u, g, conn)
}
