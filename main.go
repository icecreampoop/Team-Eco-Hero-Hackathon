package main

import (
	"log"

	backend "github.com/icecreampoop/Team-Eco-Hero-Hackathon/Backend"
)

func main() {
	
	log.Println("App is ready!")

	// listenandserve in serverhandler.go, please run this last to ensure all other functions run first
	backend.ServerHandler()
}
