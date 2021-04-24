package middle

import (
	"fmt"
	"github.com/DAT4/backend-project/dao"
)

var G *Game

func init() {
	G = NewGame()
}

type Game struct {
	state      GameState
	Db         dao.DBase
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

func (g *Game) Run(db dao.DBase) {
	g.Db = db
	for {
		select {
		case client := <-g.register:
			g.clients[client] = byte(client.Id)
		case client := <-g.unregister:
			if _, ok := g.clients[client]; ok {
				delete(g.clients, client)
				close(client.send)
			}
		case message := <-g.broadcast:
			g.sendMessageToAllClients(message)
		}
	}
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
