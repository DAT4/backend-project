// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middle

import (
	"github.com/DAT4/backend-project/models"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id   int
	user *models.User
	game *Game
	conn *websocket.Conn
	send chan []byte
}

func NewClient(g *Game, conn *websocket.Conn) {
	user, err := authenticateClient(conn, g)
	if err != nil {
		return
	}

	c := &Client{
		Id:   user.PlayerID,
		user: user,
		game: g,
		conn: conn,
		send: make(chan []byte, 256),
	}
	c.game.register <- c

	err = c.sendStartCommand(g)
	if err != nil {
		return
	}

	go c.writePump()
	go c.readPump()
}

func (c *Client) sendStartCommand(g *Game) error {
	msg := message{
		command:  CREATE,
		playerId: byte(c.Id),
		x:        1,
		y:        1,
	}

	players := make([]byte, 0, len(g.clients))

	for _, id := range g.clients {
		players = append(players, id)
	}

	return c.conn.WriteMessage(websocket.BinaryMessage, msg.sendWithContent(players))
}

func authenticateClient(c *websocket.Conn, g *Game) (u *models.User, err error) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return nil, err
		}
		if message[3] == 0 {
			token, err := getToken(string(message[3:]))
			if err != nil {
				return nil, err
			}
			u, err := UserFromToken(token, g.Db)
			if err != nil {
				c.WriteMessage(websocket.BinaryMessage, []byte{0, 0, 0, 5, 1})
				return nil, err
			}
			return &u, nil
		}
	}
}

func getToken(token string) (string, error) {
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) < 2 {
		return tokenParts[0], nil
	}
	return tokenParts[1], nil
}

func (c *Client) readPump() {
	defer func() {
		c.game.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if message[3] == 0 {
			u, err := UserFromToken(string(message[3:]), c.game.Db)
			if err != nil {
				return
			}
			c.user = &u
			continue
		}
		c.game.broadcast <- message //[id][x][y]//[command][message/string]
	}

}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.BinaryMessage, message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}
