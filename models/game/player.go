package game

import (
	"github.com/DAT4/backend-project/models/user"
)

type Player struct {
	*user.User
	Position
}

type Position struct {
	x float64
	y float64
}

