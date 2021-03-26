package game

import (
	"fmt"
)

type Game struct {
	State      GameState
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
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
			g.Clients[client] = true
			fmt.Println("Player successfully connected")
		case client := <-g.Unregister:
			if _, ok := g.Clients[client]; ok {
				delete(g.Clients, client)
				close(client.Send)
			}
		case message := <-g.Broadcast:
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
	for client := range g.Clients {
		select {
		case client.Send <- message:
			fmt.Println("Sending to client")
		default:
			close(client.Send)
			delete(g.Clients, client)
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
