package main

import (
	"fmt"
	"texting-app/internal/app"
)

func main() {
	error := app.Run()
	if error != nil {
		fmt.Println("App module unable to start.")
	} else {
		fmt.Println("App module started successfully.")
	}
}
