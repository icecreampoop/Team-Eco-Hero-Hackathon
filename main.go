package main

import (
	"log"

	backend "github.com/icecreampoop/Team-Eco-Hero-Hackathon/Backend"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	//backend.LoadDataFromConfig("data.json")
	//log.Printf("Loaded users: %+v", backend.Users)
	//log.Printf("Loaded items: %+v", backend.Items)
	//backend.AddNewItem(1, "SUPER Cable", "electronics")
	//backend.DeleteItem(6)
	//backend.EditItem(5, "COOL JACKET", "", "abcd", "")
	log.Println("App is ready!")

	// listenandserve in serverhandler.go, please run this last to ensure all other functions run first
	backend.ServerHandler()
}
