package main

import (
	"flag"
	"github.com/levgenia/websocket-chat/admin"
	"github.com/levgenia/websocket-chat/chat"
	"os"
	"strconv"
)

var (
	port    int
)

func init() {
	flag.IntVar(&port, "port", 3000, "host port")
}

func main() {
	flag.Parse()
	os.Setenv("PORT", strconv.Itoa(port))

	clients := chat.NewClients()
	go chat.Server(port, clients)

	admin := admin.NewAdmin(clients)
	admin.ListenConsole()

}
