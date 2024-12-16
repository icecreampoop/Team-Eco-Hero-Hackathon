package main

import (
	backend "github.com/icecreampoop/Team-Eco-Hero-Hackathon/Backend"
	"github.com/joho/godotenv"
)

import (
	_ "github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	//backend.ServerHandler()
	backend.LoadDataFromConfig("data.json")
	log.Printf("Loaded users: %+v", backend.Users)
	log.Printf("Loaded items: %+v", backend.Items)

	log.Println("App is ready!")
}
