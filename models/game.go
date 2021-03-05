package models

import "fmt"

type GameState int8

const (
	Opening GameState = iota
	Full
	Empty
	Closing
)

type Game struct {
	Name       Username
	state      GameState
	clients    map[*Client]bool
	players    [2]Player
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewGame(name Username) *Game {
	return &Game{
		Name:       name,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (g *Game) Run() {
	for {
		select {
		case client := <-g.register:
			g.clients[client] = true

		case client := <-g.unregister:
			if _, ok := g.clients[client]; ok {
				delete(g.clients, client)
				close(client.send)
			}
		case message := <-g.broadcast:
			g.Parse(message)
		}
	}
}

func (g *Game) Parse(msg []byte) {
	//TODO use this to get commands in byte
	//Maybe this should be done in the client if only related to him
	//But if it is related to everyone then here
	g.sendMessageToAllClients(msg)
}

func (g *Game) sendMessageToAllClients(message []byte) {
	for client := range g.clients {
		select {
		case client.send <- message:
			fmt.Println("Sending to client")
		default:
			close(client.send)
			delete(g.clients, client)
		}
	}
}
