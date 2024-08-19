package app

import (
	"net/http"
)

func mapRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	miscController := MiscController{}
	wsController := WSController{}

	mux.HandleFunc("/test", miscController.test)
	mux.HandleFunc("/ws", wsController.handleConn)

    go wsController.handleMessages()

	return mux
}

func Run() error {
	router := mapRoutes()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()

	return nil
}
