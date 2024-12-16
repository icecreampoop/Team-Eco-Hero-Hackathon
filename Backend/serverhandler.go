package backend

import (
	"fmt"
	"log"
	"net/http"
)

// functions to handle HTTP requests for page loads
func HandleHTTPIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/index.gohtml")
}

func HandleHTTPItems(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/items.gohtml")
}

func HandleHTTPUser(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/user.gohtml")
}

func HandleHTTPBoard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/board.gohtml")
}

func HandleHTTPLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/login.html")
}

func HandleHTTPSignup(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/signup.html")
}

// // functions to handle user actions
// func HandleButtonClick(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		fmt.Fprintf(w, "Button clicked!")
// 	} else {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }

// func HandleUserInput(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		// Parse form data
// 		if err := r.ParseForm(); err != nil {
// 			http.Error(w, "Unable to parse form", http.StatusBadRequest)
// 			return
// 		}
// 		// Get user input
// 		userInput := r.FormValue("userInput")
// 		fmt.Fprintf(w, "Received user input: %s", userInput)
// 	} else {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }

func ServerHandler() {
	// Create new HTTP mux
	mux := http.NewServeMux()

	// handlers
	mux.HandleFunc("/", HandleHTTPIndex) // default handler to index
	mux.HandleFunc("/index", HandleHTTPIndex)
	mux.HandleFunc("/items", HandleHTTPItems)
	mux.HandleFunc("/user", HandleHTTPUser)
	mux.HandleFunc("/board", HandleHTTPBoard)
	mux.HandleFunc("/login", HandleHTTPLogin)
	mux.HandleFunc("/signup", HandleHTTPSignup)

	// Serve static files from the frontend directory
	fs := http.FileServer(http.Dir("./Frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
