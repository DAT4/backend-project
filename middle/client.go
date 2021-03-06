// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middle

import (
	"fmt"
	"github.com/DAT4/backend-project/models"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id   byte
	user *models.User
	game *Game
	conn *websocket.Conn
	send chan []byte
}

func NewClient(g *Game, conn *websocket.Conn) {
	user, err := authenticateClient(conn, g)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k := range g.clients {
		_ = conn.WriteMessage(websocket.BinaryMessage, message{
			command:  ASSIGN,
			playerId: byte(k.user.PlayerID),
			x:        0,
			y:        0,
		}.send())
	}

	c := &Client{
		Id:   byte(user.PlayerID),
		user: user,
		game: g,
		conn: conn,
		send: make(chan []byte, 256),
	}
	c.game.register <- c

	err = c.sendStartCommand()
	if err != nil {
		return
	}

	go c.writePump()
	go c.readPump()
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
		switch message[ACT] {
		case READY:
			u, err := UserFromToken(string(message[ACT:]), c.game.Db)
			if err != nil {
				return
			}
			c.user = &u
			continue
		case MOVE:
			if c.game.onTheRoad(message) {
				c.game.broadcast <- message //----[id][x][y][command][message/string]
			}

		default:
			fmt.Println("ERROR:", message)
		}
	}

}
func (c *Client) sendStartCommand() error {
	msg := message{
		command:  CREATE,
		playerId: c.Id,
		x:        1,
		y:        1,
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, msg.send())
}

func (g *Game) onTheRoad(msg []byte) (ok bool) {
	x, y := int(msg[X]), int(msg[Y])
	if x > -1 && x < 30 && y > -1 && y < 30 {
		return 0 != g.Map[1][x+y*30]
	} else {
		return false
	}
}
func authenticateClient(c *websocket.Conn, g *Game) (u *models.User, err error) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return nil, err
		}
		if message[ACT] == READY {
			token, err := getToken(string(message[ACT:]))
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
