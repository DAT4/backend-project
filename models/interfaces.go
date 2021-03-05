package models

type Router interface {
	Run()
	unregister(*Client)
	register(*Client)
	broadcast([]byte)
}
