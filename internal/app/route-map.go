package app

import (
	"net/http"
	"texting-app/internal/api/handlers"
	"texting-app/internal/api/middlewares"
)

func MapRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	wsController := handlers.WSController{}
	pubFilesHandler := http.FileServer(http.Dir("./static"))

	mux.Handle("/public/",http.StripPrefix("/public/", pubFilesHandler))
	mux.HandleFunc("/ws", wsController.HandleConn)
	mux.HandleFunc("/login", handlers.Login)
    mux.HandleFunc("/signup", handlers.Signup)
	mux.HandleFunc("/chats", handlers.Chats)
    mux.HandleFunc("/hx/chat-msg-list", middlewares.HTMXOnly(handlers.ChatMsgList))

	go wsController.HandleMessages()

	return mux
}
