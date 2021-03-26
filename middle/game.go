package middle

import (
	"github.com/DAT4/backend-project/models/game"
)

var Game *game.Game

func init() {
	Game = game.NewGame()
}
