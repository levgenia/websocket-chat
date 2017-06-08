package main

import (
	"github.com/codegangsta/martini"
	"github.com/levgenia/websocket-chat/chat"
	"bufio"
	"os"
	"log"
	"github.com/levgenia/websocket-chat/admin"
	"strconv"
)

const (
	portNum = 3000
)

func init() {
	os.Setenv("PORT", strconv.Itoa(portNum))
}

func main() {
	m := martini.Classic()

	f, err := os.OpenFile("logfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	logger := log.New(f, "", 1)
	m.Logger(logger)

	m.Get("/", func() string {
		return `<html><body><script src='//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js'></script>
    <ul id=messages></ul><form><input id=message><input type="submit" id=send value=Send></form>
    <script>
    var c=new WebSocket('ws://localhost:` + strconv.Itoa(portNum) + `/sock');
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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		out := admin.Handle(scanner.Text())
		os.Stdout.Write([]byte(out))
	}
}
