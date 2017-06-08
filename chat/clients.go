package chat

import (
	"sync"
)

type Clients struct {
	sync.RWMutex
	connections map[*connection]int
}

func NewClients() *Clients {
	return &Clients{
		connections: make(map[*connection]int),
	}
}

func (cc *Clients) AddConnection(conn *connection) {
	cc.Lock()
	cc.connections[conn] = 0
	cc.Unlock()
}

func (cc *Clients) DeleteConnection(conn *connection) {
	cc.Lock()
	delete(cc.connections, conn)
	cc.Unlock()
}

func (cc *Clients) DeleteByAddress(address string) bool {
	cc.Lock()
	defer cc.Unlock()
	for k := range cc.connections {
		if k.address.String() == address {
			delete(cc.connections, k)
			return true
		}
	}
	return false
}

func (cc *Clients) ExistConnection(conn *connection) bool {
	cc.RLock()
	defer cc.RUnlock()

	_, ok := cc.connections[conn]
	return ok
}

func (cc *Clients) ListAll() (result []string) {
	for k := range cc.connections {
		result = append(result, k.address.String())
	}
	return
}

func (cc *Clients) BroadcastMessage(messageType int, message []byte) {
	cc.RLock()
	defer cc.RUnlock()

	for client := range cc.connections {
		if err := client.SendMessage(messageType, message); err != nil {
			return
		}
	}
}
