package game

import (
	"fmt"
	"github.com/DAT4/backend-project/models/user"
)

type Game struct {
	Name       user.Username
	state      GameState
	clients    map[*Client]bool
	players    [2]*Player
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
}

func NewGame() *Game {
	return &Game{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func assignPlayer(number int) message {
	return message{
		command:  ASSIGN,
		playerId: byte(number),
		startPos: Position{
			x: float64((number + 1) * 40),
			y: float64((number + 1) * 40),
		},
	}
}

func (g *Game) checkReady() bool {
	var ready = true
	for i := 0; i < 2; i++ {
		if g.players[i] == nil {
			ready = false
		}
	}
	return ready
}

func (g *Game) Run() {
	for {
		select {
		case client := <-g.Register:
			fmt.Println("Trying to connect new player")
			err := client.SendStartCommand(assignPlayer(client.User.PlayerID))
			if err != nil {
				fmt.Println("error with json closing ws:", err)
				close(client.Send)
				return
			}
			g.clients[client] = true
			fmt.Println("Player successfully connected")
		case client := <-g.unregister:
			if _, ok := g.clients[client]; ok {
				delete(g.clients, client)
				close(client.Send)
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
		case client.Send <- message:
			fmt.Println("Sending to client")
		default:
			close(client.Send)
			delete(g.clients, client)
		}
	}
}
