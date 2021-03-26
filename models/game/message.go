package game

type message struct {
	command  byte
	playerId byte
	x        byte
	y        byte
}

func (m message) sendWithContent(content string) (out []byte) {
	return append(m.send(), []byte(content)...)
}

func (m message) send() []byte {
	return []byte{m.playerId, m.x, m.y, m.command}
}
