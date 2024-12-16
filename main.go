package main

import backend "github.com/icecreampoop/Team-Eco-Hero-Hackathon/Backend"

func main() {
	backend.UploadFile("test.txt", []byte("test"))
}

