// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

import (
	"github.com/DAT4/backend-project/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	user *models.User
	game *Game
	conn *websocket.Conn
	send chan []byte
}

func NewClient(u *models.User, g *Game, conn *websocket.Conn) {
	c := &Client{}
	c.init(u, g, conn)
	go c.writePump()
	go c.readPump()
}

func (c *Client) init(u *models.User, g *Game, conn *websocket.Conn) {
	c.user = u
	c.game = g
	c.conn = conn
	c.send = make(chan []byte, 256)
	c.game.register <- c
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

func (c *Client) sendStartCommand(msg message, others []byte) error {
	err := c.conn.WriteMessage(websocket.BinaryMessage, msg.sendWithContent(others))
	if err != nil {
		return err
	}
	return nil
}
