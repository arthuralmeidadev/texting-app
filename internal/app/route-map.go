package app

import (
	"net/http"
	"texting-app/internal/app/handlers"
)

func MapRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	wsController := handlers.WSController{}
	pubFilesHandler := http.FileServer(http.Dir("./static"))

	mux.Handle("/public/", http.StripPrefix("/public/", pubFilesHandler))
	mux.HandleFunc("/ws", EnsureAuth(wsController.HandleConn))
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/signup", handlers.Signup)
	mux.HandleFunc("/chats", EnsureAuth(handlers.Chats))
	mux.HandleFunc("/hx/chat-msg-list", EnsureAuth(HTMXOnly(handlers.ChatMsgList)))
	mux.HandleFunc("/hx/new-chat", EnsureAuth(HTMXOnly(handlers.NewChat)))
	mux.HandleFunc("/hx/friend-list", EnsureAuth(HTMXOnly(handlers.FriendList)))
	mux.HandleFunc("/hx/user-list", EnsureAuth(HTMXOnly(handlers.UserList)))

	go wsController.HandleMessages()

	return mux
}
