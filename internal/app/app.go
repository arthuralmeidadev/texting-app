package app

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func Run() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	router := MapRoutes()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

    err = server.ListenAndServe()
    if err != nil {
        return err
    }

	return nil
}
