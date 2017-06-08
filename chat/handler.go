package chat

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

func ChatHandler(clients *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a web socket handshake", 400)
			return
		} else if err != nil {
			log.Println(err)
			return
		}
		client := ws.RemoteAddr()
		sockCli := NewConnection(ws, client)

		for {
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				log.Println("failed to read meddage", client.String(), err)
				clients.DeleteConnection(sockCli)
				return
			}

			if !clients.ExistConnection(sockCli) {
				if IsStartMessage(p) {
					clients.AddConnection(sockCli)
					sockCli.SendMessage(messageType, []byte("hello, you're in chat! type '/stop' to stop chatting"))
				} else {
					sockCli.SendError(messageType, []byte("type '/start' to start chatting"))
				}
			} else {
				if IsStopMessage(p) {
					clients.DeleteConnection(sockCli)
					sockCli.SendMessage(messageType, []byte("bye-bye"))
				} else {
					clients.BroadcastMessage(messageType, p)
				}
			}
		}
	}
}

