package api

import (
	"fmt"
	"github.com/DAT4/backend-project/middle"
	"net/http"
)

func joinWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Peer")
	u, err := middle.UserFromToken(r)
	if err != nil {
		handleHttpError(w, err, http.StatusNotAcceptable)
	}
	middle.ServeWs(&u, middle.Game, w, r)
}
