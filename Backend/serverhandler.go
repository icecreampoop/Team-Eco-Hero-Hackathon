package backend

import (
	"fmt"
	"log"
	"net/http"
)

func HandleHTTPIndex(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Index")
	http.ServeFile(w, r, "frontend/index.gohtml")
}

func HandleHTTPItems(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Items")
	// http.ServeFile(w, r, "frontend/items.gohtml")
}

func HandleHTTPUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User")
	// http.ServeFile(w, r, "frontend/user.gohtml")
}

func HandleHTTPBoard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Leaderboard")
	// http.ServeFile(w, r, "frontend/board.gohtml")
}

func HandleHTTPLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login")
	// http.ServeFile(w, r, "frontend/login.gohtml")
}

func ServerHandler() {
	// Create new HTTP mux
	mux := http.NewServeMux()

	// handlers
	mux.HandleFunc("/", HandleHTTPIndex) // default handler
	mux.HandleFunc("/index", HandleHTTPIndex)
	mux.HandleFunc("/items", HandleHTTPItems)
	mux.HandleFunc("/user", HandleHTTPUser)
	mux.HandleFunc("/board", HandleHTTPBoard)
	mux.HandleFunc("/login", HandleHTTPLogin)

	// Serve static files from the frontend directory
	fs := http.FileServer(http.Dir("./frontend")) // default relative directory
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	// Start server
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
