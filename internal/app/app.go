package app

import (
	"net/http"
)

func Run() error {
	router := MapRoutes()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()

	return nil
}
