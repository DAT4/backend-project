package api

import (
	"github.com/DAT4/backend-project/middle"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *API) JoinWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	serveWs(a.Game, w, r)
}

func serveWs(g *middle.Game, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	middle.NewClient(g, conn)
}
