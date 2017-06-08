package admin

import (
	"github.com/levgenia/websocket-chat/chat"
	"fmt"
)

type admin struct {
	clients *chat.Clients
}

func NewAdmin(clients *chat.Clients) *admin {
	return &admin{
		clients: clients,
	}
}

func (a *admin) Handle(in string) (out string) {
	if in == "/list" {
		for k, v := range a.clients.ListAll() {
			out += fmt.Sprintln("User", k, ":", v)
		}
		return

	} else {
		deleted := a.clients.DeleteByAddress(in)
		if deleted {
			return fmt.Sprintln("deleted", in, "success")
		} else {
			return fmt.Sprintln("couldn't delete", in)
		}
	}
}
