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

func (ctl *WSController) handleConn(w http.ResponseWriter, r *http.Request) {
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
			ctl.clients[msg.Username] = ws
			usrIsSet = true
		}

		ctl.broadcast <- msg
	}
}

func (ctl *WSController) handleMessages() {
	msg := <-ctl.broadcast
	for _, recip := range msg.Recipients {
		cli := ctl.clients[recip]
        if cli == nil {
            continue
        }

		err := cli.ReadJSON(msg.Message)
		if err != nil {
			log.Printf("error: %v", err)
			cli.Close()
			delete(ctl.clients, recip)
		}
	}
}
