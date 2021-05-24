package middle

import "time"

type TokenType int

const (
	ID = iota
	X
	Y
	ACT
	DIRECTION
)

const (
	AUTHENTICATION TokenType = iota
	REFRESH
)

type GameState int8

const (
	Opening GameState = iota
	Full
	Empty
	Closing
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

const (
	READY byte = iota
	CREATE
	ASSIGN
	MOVE
	WRITE
)

//Directions
const (
	LEFT = iota
	RIGHT
	UP
	DOWN
)
