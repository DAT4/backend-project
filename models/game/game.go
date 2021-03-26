package game

import (
	"fmt"
)

type Game struct {
	state      GameState
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewGame() *Game {
	return &Game{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (g *Game) Run() {
	for {
		select {
		case client := <-g.register:
			fmt.Println("Trying to connect new player")
			err := client.sendStartCommand(assignPlayer(client.user.PlayerID))
			if err != nil {
				fmt.Println("error with json closing ws:", err)
				close(client.send)
				return
			}
			g.clients[client] = true
			fmt.Println("Player successfully connected")
		case client := <-g.unregister:
			if _, ok := g.clients[client]; ok {
				delete(g.clients, client)
				close(client.send)
			}
		case message := <-g.broadcast:
			g.parse(message)
		}
	}
}

func (g *Game) parse(msg []byte) {
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

func assignPlayer(number int) message {
	return message{
		command:  ASSIGN,
		playerId: byte(number),
		x:        0,
		y:        0,
	}
}
