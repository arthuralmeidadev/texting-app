package main

import (
	"log"
	"texting-app/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
        log.Println("App module unable to start. Reason:")
        log.Println(err)
	} else {
		log.Println("App module started successfully.")
	}
}
