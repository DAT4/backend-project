package game

import (
	"fmt"
)

type Game struct {
	state      GameState
	counter    int
	clients    map[*Client]byte
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewGame() *Game {
	return &Game{
		clients:    make(map[*Client]byte),
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
			g.counter++
			players := make([]byte, 0, len(g.clients))
			for _, id := range g.clients{
				players = append(players, id)
			}
			err := client.sendStartCommand(assignPlayer(g.counter), players)
			if err != nil {
				fmt.Println("error with json closing ws:", err)
				close(client.send)
				return
			}
			g.clients[client] = byte(g.counter)
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
	fmt.Println(message)

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
		command:  CREATE,
		playerId: byte(number),
		x:        0,
		y:        0,
	}
}
