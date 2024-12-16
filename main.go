package main

import (
	backend "github.com/icecreampoop/Team-Eco-Hero-Hackathon/Backend"
	// "github.com/joho/godotenv"
	// "log"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }
	backend.ServerHandler()
	//backend.Connect()
}
