package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func ChatHandler(clients *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a web socket handshake", 400)
			return
		} else if err != nil {
			// TODO: wrap err with more informative message
			log.Println(err)
			return
		}
		client := ws.RemoteAddr()
		clientConn := NewConnection(ws, client)

		for {
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				log.Println("failed to read message", client.String(), err)
				clients.DeleteConnection(clientConn)
				return
			}

			if !clients.ExistConnection(clientConn) {
				if IsStartMessage(p) {
					str := fmt.Sprint(clientConn.address.String(), " joined to the chat")
					clients.BroadcastMessage(messageType, []byte(str))
					clients.AddConnection(clientConn)
					clientConn.SendMessage(messageType, []byte("hello, you're in chat! type '/stop' to stop chatting"))
				} else {
					clientConn.SendError(messageType, []byte("type '/start' to start chatting"))
				}
			} else {
				if IsStopMessage(p) {
					clients.DeleteConnection(clientConn)
					clientConn.SendMessage(messageType, []byte("bye-bye! I'll miss u :("))
					str := fmt.Sprint(clientConn.address.String(), " left the chat")
					clients.BroadcastMessage(messageType, []byte(str))
				} else {
					clients.BroadcastMessage(messageType, p)
				}
			}
		}
	}
}
