package models

type Player struct {
	*User
	Position
}

type Position struct {
	x float64
	y float64
}

