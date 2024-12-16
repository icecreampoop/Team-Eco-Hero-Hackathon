package main

import (
	"fmt"
	"log"
	"net/http"
)

func HandleHTTPIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index")
}

func HandleHTTPItems(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Items")
}

func HandleHTTPUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User")
}

func HandleHTTPBoard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Leaderboard")
}

func HandleHTTPLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login")
}

func main() {
	// Create new HTTP mux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/", HandleHTTPIndex) // default handler
	mux.HandleFunc("/index", HandleHTTPIndex)
	mux.HandleFunc("/items", HandleHTTPItems)
	mux.HandleFunc("/user", HandleHTTPUser)
	mux.HandleFunc("/board", HandleHTTPBoard)
	mux.HandleFunc("/login", HandleHTTPLogin)

	// Start server
	port := ":8080"
	fmt.Printf("Server started at port %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
