package backend

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
)

func showAllItems(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/items.html")
}

func showSingleItem(w http.ResponseWriter, r *http.Request) {

}

func createNewItem(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Retrieve the file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Validate file type (optional, based on MIME type)
	contentType := handler.Header.Get("Content-Type")
	if contentType[:6] != "image/" {
		http.Error(w, "Only image files are allowed", http.StatusUnsupportedMediaType)
		return
	}

	// Read file into []byte
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	imageBytes := buf.Bytes()

	// Detect the image format
	_, format, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Unsupported or invalid image format", http.StatusUnsupportedMediaType)
		return
	}

	// Process the imageBytes (e.g., store in a database or perform operations)
	fmt.Printf("Received file %s with size %d bytes\n", handler.Filename, len(imageBytes))

	// UploadFile(itemName + "."  + format, imageBytes)
	UploadFile("."+format, imageBytes)

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded and processed successfully"))
}

func requestItem(w http.ResponseWriter, r *http.Request) {

}

func acceptRequest(w http.ResponseWriter, r *http.Request) {

}

func updateItemDetails(w http.ResponseWriter, r *http.Request) {

}

func deleteItem(w http.ResponseWriter, r *http.Request) {

}

// functions to handle HTTP requests for page loads
func HandleHTTPIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/index.gohtml")
}

func HandleHTTPUser(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/user.html")
}

func HandleHTTPBoard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/board.html")
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
	// mux.HandleFunc("/", HandleHTTPIndex) // default handler to index
	// mux.HandleFunc("/index", HandleHTTPIndex)
	// Matt: Changed the / path to showallitems instead
	mux.HandleFunc("/", showAllItems)

	//all item handlers
	mux.HandleFunc("GET /items", showAllItems)
	mux.HandleFunc("GET /items/{itemID}", showSingleItem)
	mux.HandleFunc("POST /items", createNewItem)
	mux.HandleFunc("POST /items/{itemID}/request", requestItem)
	mux.HandleFunc("POST /items/{itemID}/accept", acceptRequest)
	mux.HandleFunc("PUT /items/{itemID}", updateItemDetails)
	mux.HandleFunc("DELETE /items/{itemID}", deleteItem)

	mux.HandleFunc("/user", HandleHTTPUser)
	mux.HandleFunc("/board", HandleHTTPBoard)
	mux.HandleFunc("/login", HandleHTTPLogin)
	mux.HandleFunc("/signup", HandleHTTPSignup)

	// Serve static files from the frontend directory
	fs := http.FileServer(http.Dir("./Frontend/static"))
	mux.Handle("/Frontend/static/", http.StripPrefix("/Frontend/static/", fs))

	// Start server
	port := ":5000"
	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
