package main

import "fmt"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	spaces     chan bool
	roles      []*Client
}

func newHub() *Hub {
	h := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		roles:      make([]*Client, 2),
		spaces:     make(chan bool, 2),
	}
	h.spaces <- true
	h.spaces <- true
	fmt.Println("Hub for 2 created")
	return h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					fmt.Println("Sending to client")
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
