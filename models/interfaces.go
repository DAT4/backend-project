package models

import "github.com/DAT4/backend-project/models/game"

type Router interface {
	Run()
	unregister(*game.Client)
	register(*game.Client)
	broadcast([]byte)
}
