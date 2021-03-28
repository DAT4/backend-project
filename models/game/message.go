package game

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
	return []byte{m.playerId, m.x, m.y, m.command}
}
