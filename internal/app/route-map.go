package app

import (
	"net/http"
	"texting-app/internal/handlers"
)

func MapRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	wsController := handlers.WSController{}

	mux.HandleFunc("/ws", wsController.HandleConn)
	mux.HandleFunc("/login", handlers.Login)

	go wsController.HandleMessages()

	return mux
}
