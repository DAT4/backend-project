package middle

import "fmt"

type message struct {
	command  byte
	playerId byte
	x        byte
	y        byte
}

func (m message) sendWithContent(command []byte) (out []byte) {
	return append(m.send(), command...)
}

func (m message) send() []byte {
	msg := []byte{m.playerId, m.x, m.y, m.command}
	fmt.Println(msg)
	return msg
}
