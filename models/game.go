package models

import (
	"fmt"
)

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
	players    [2]*Player
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewGame(name Username) *Game {
	return &Game{
		Name:       name,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

const (
	READY = iota
	CREATE
	ASSIGN
)

type message struct {
	command  byte
	playerId byte
	startPos Position
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
		case client := <-g.register:
			fmt.Println("Trying to connect new player")
			if x, ok := g.checkAndAppend(client.player); ok {
				err := client.SendStartCommand(assignPlayer(x))
				if err != nil {
					fmt.Println("error with json closing ws:", err)
					close(client.send)
					return
				}
				g.clients[client] = true
				fmt.Println("Player successfully connected")
			} else {
				close(client.send)
			}
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

func (g *Game) checkAndAppend(p *Player) (int, bool) {
	var space bool
	var x int
	for i := 0; i < 2; i++ {
		if g.players[i] == nil {
			x = i
			space = true
		} else if p.Username == g.players[i].Username {
			return i, true
		}
	}
	if space {
		g.players[x] = p
		return x, true
	}
	return 0, false
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
