package handlers

import (
	"log"
	"net/http"
	"texting-app/internal/pkg/providers"
	"texting-app/internal/pkg/store"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type msgDTO struct {
	Sender    string
	Content   string `json:"value"`
	ChatId    uint   `json:"recipient"`
	RepliesTo int    `json:"repliesTo"`
}

type WSController struct {
	clients map[string]*websocket.Conn
	msgChan chan msgDTO
}

func (ctl *WSController) HandleConn(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	usrIsSet := false
	for {
		var msgDto msgDTO
		err := ws.ReadJSON(&msgDto)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		usrn := r.Header.Get("username")
		msgDto.Sender = usrn

		if !usrIsSet {
			ctl.clients[usrn] = ws
			usrIsSet = true
		}

		ctl.msgChan <- msgDto
	}
}

func (ctl *WSController) HandleMessages() {
	storeInst, err := store.GetStore()
	if err != nil {
		panic(err.Error())
	}
	msgProvider := providers.NewMessageProvider(storeInst)
	chatProvider := providers.NewChatProvider(storeInst)
	for {
		msgDto := <-ctl.msgChan
		msgStoreChan := make(chan struct{})
		membersChan := make(chan []string)
		go func() {
			msgProvider.StoreMessage(
				msgDto.Sender,
				msgDto.Content,
				msgDto.ChatId,
				msgDto.RepliesTo,
			)
			msgStoreChan <- struct{}{}
		}()

		go func() {
			members, _ := chatProvider.GetChatMembers(msgDto.ChatId)
			membersChan <- members
		}()

		for _, member := range <-membersChan {
			if member == msgDto.Sender {
				continue
			}

			cli := ctl.clients[member]
			if cli == nil {
				continue
			}

			err := cli.ReadJSON(msgDto.Content)
			if err != nil {
				log.Printf("error: %v", err)
				cli.Close()
			}
		}
        <-msgStoreChan
	}
}
