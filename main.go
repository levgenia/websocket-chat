package main

import (
	"flag"
	"github.com/codegangsta/martini"
	"github.com/levgenia/websocket-chat/admin"
	"github.com/levgenia/websocket-chat/chat"
	"log"
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

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return
	}
	defer f.Close()

	m := martini.Classic()

	logger := log.New(f, "", 1)
	m.Logger(logger)

	m.Get("/", func() string {
		return `<html><body><script src='//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js'></script>
    <ul id=messages></ul>
    <form> <input id=message><input type="submit" id=send value=Send></form>
    <script>
    var addr = 'ws://' + window.location.host + '/sock'
    var c=new WebSocket(addr);
    c.onopen = function(){
      c.onmessage = function(response){
        console.log(response.data);
        var newMessage = $('<li>').text(response.data);
        $('#messages').append(newMessage);
        $('#message').val('');
      };
      $('form').submit(function(){
        c.send($('#message').val());
        return false;
      });
    }
    </script></body></html>`
	})

	clients := chat.NewClients()
	admin := admin.NewAdmin(clients)

	m.Get("/sock", chat.ChatHandler(clients))

	go m.Run()

	admin.ListenConsole()

}
