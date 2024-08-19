package app

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type reqBody struct {
	Username   string   `json:"username"`
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
}

type WSController struct {
	clients   map[string]*websocket.Conn
	broadcast chan reqBody
}

func (self *WSController) handleConn(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	usrIsSet := false
	for {
		var msg reqBody
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		if !usrIsSet {
			self.clients[msg.Username] = ws
			usrIsSet = true
		}

		self.broadcast <- msg
	}
}

func (self *WSController) handleMessages() {
	msg := <-self.broadcast
	for _, recip := range msg.Recipients {
		client := self.clients[recip]
		err := client.ReadJSON(msg.Message)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(self.clients, recip)
		}
	}
}
