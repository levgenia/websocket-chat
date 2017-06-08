package chat

import (
	"net"
	"github.com/gorilla/websocket"
)

type connection struct {
	websocket *websocket.Conn
	address   net.Addr
}

func NewConnection(ws *websocket.Conn, addr net.Addr) *connection {
	return &connection{
		ws,
		addr,
	}
}

func (c *connection) SendMessage(messageType int, message []byte) error {
	out := append([]byte(c.address.String() + ": "), message...)
	return c.websocket.WriteMessage(messageType, out)
}

func (c *connection) SendError(messageType int, message []byte) error {
	return c.websocket.WriteMessage(messageType, message)
}
