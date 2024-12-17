package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

var tpl *template.Template
var hasher = fnv.New32a()

type userPageData struct {
	TplUser  User
	TplItems []Item
}

func showAllItems(w http.ResponseWriter, r *http.Request) {
	// Load data from data.json
	data, err := LoadUserData()
	if err != nil {
		http.Error(w, "Error loading data", http.StatusInternalServerError)
		log.Println("Error loading data:", err)
		return
	}

	// Get the filter query parameters from the URL (if any)
	statusFilter := r.URL.Query().Get("status") // "available" or "pending"

	// Filter items based on the status
	var filteredItems []Item
	for _, item := range data.Items {
		// Filter by status if specified, otherwise show all items
		if statusFilter == "" || item.ItemStatus == StatusAvailable || item.ItemStatus == StatusPending {
			filteredItems = append(filteredItems, item)
		}
	}

	// Pass the filtered items to the template
	err = tpl.ExecuteTemplate(w, "items.html", filteredItems)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func showSingleItem(w http.ResponseWriter, r *http.Request) {
	asf
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

	// upload media to digital ocean spaces
	// Get the "UserID" cookie from the request
	cookie, err := r.Cookie("UserID")
	if err != nil {
		// If the cookie is not found, handle the error
		http.Error(w, "UserID cookie not found", http.StatusBadRequest)
		return
	}

	// Convert the cookie value (which is a string) to an integer
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		// If there's an error converting the value, handle it
		http.Error(w, "Invalid UserID value", http.StatusBadRequest)
		return
	}
	hashedFileName := hashResourcePath(findUser(userID).Email+r.FormValue("item-name")) + "." + format
	fileResourcePath, _ := UploadFile(hashedFileName, imageBytes)
	// add item entry to db
	userIDInt, _ := getUserID(r)
	AddNewItem(userIDInt, r.FormValue("item-name"), r.FormValue("item-description"),
		r.FormValue("category"), fileResourcePath)

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded and processed successfully"))
}

func requestItem(w http.ResponseWriter, r *http.Request) {

}

func acceptRequest(w http.ResponseWriter, r *http.Request) {

}

func updateItemDetails(w http.ResponseWriter, r *http.Request) {
	asf
}

func deleteItem(w http.ResponseWriter, r *http.Request) {

}

// functions to handle HTTP requests for page loads
func HandleHTTPIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/index.gohtml")
}

func HandleHTTPUser(w http.ResponseWriter, r *http.Request) {
	// redirect function

	// userID, err := getUserID(r)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	userID := 2

	data, err := LoadUserData()
	if err != nil {
		fmt.Println(err)
		return
	}

	var foundUser User
	for _, user := range data.Users {
		if user.UserID == userID {
			foundUser = user
		}
	}

	// Get user's listings
	var userItems []Item
	for _, item := range data.Items {
		if item.OwnerID == userID {
			userItems = append(userItems, item)
		}
	}
	tplData := userPageData{
		foundUser,
		userItems,
	}

	err = tpl.ExecuteTemplate(w, "user.html", tplData)
	if err != nil {
		http.Error(w, "Error rendering User template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func HandleHTTPBoard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/board.html")
}

// // HandleHTTPLogin serves the login page
// func HandleHTTPLogin(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./Frontend/static/login.html")
// }

// HandleHTTPLogin serves the login page and handles login authentication
func HandleHTTPLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./Frontend/static/login.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get form values
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate credentials
		valid, err := ValidateUserCredentials(username, password)
		if err != nil {
			http.Error(w, "Unable to validate user credentials", http.StatusInternalServerError)
			return
		}

		if valid {
			// Successful login
			http.Redirect(w, r, "/user", http.StatusSeeOther)
			return
		}

		// Invalid credentials
		tmpl.Execute(w, map[string]interface{}{
			"ErrorMessage": "Invalid username or password",
		})
		return
	}

	// Serve the login page for GET requests
	tmpl.Execute(w, nil)
}

// HandleHTTPSignup serves the signup page and handles user registration
func HandleHTTPSignup(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./Frontend/static/signup.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get form values
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Add the new user to the data.json file
		err = AddNewUser(email, password)
		if err != nil {
			http.Error(w, "Unable to add new user", http.StatusInternalServerError)
			return
		}

		// Redirect to login page after successful signup
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Serve the signup page for GET requests
	tmpl.Execute(w, nil)
}

func ServerHandler() {
	// Go templates
	var err error
	tpl, err = template.ParseGlob("./Frontend/static/*.html")
	if err != nil {
		log.Println("Error parsing template:", err)
	}

	// Create new HTTP mux
	mux := http.NewServeMux()

	// Default handler
	mux.HandleFunc("/", showAllItems) // default handler to showallitems

	//all item handlers
	mux.HandleFunc("GET /items", showAllItems)
	mux.HandleFunc("GET /items/{itemID}", showSingleItem)
	mux.HandleFunc("POST /items", createNewItem)
	mux.HandleFunc("POST /items/{itemID}/request", requestItem)
	mux.HandleFunc("POST /items/{itemID}/accept", acceptRequest)
	mux.HandleFunc("POST /items/{itemID}/update-item", updateItemDetails)
	mux.HandleFunc("GET /items/{itemID}/update-item", showSingleItem)
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

func getUserID(r *http.Request) (int, error) {
	// Retrieve userID from cookie
	cookie, err := r.Cookie("userID")
	if err == http.ErrNoCookie {
		return -1, fmt.Errorf("no userID cookie found. Please log in")
	} else if err != nil {
		return -1, fmt.Errorf("error retrieving cookie")
	}

	userID, _ := strconv.Atoi(cookie.Value)

	return userID, nil
}

func loadUsers(filename string) ([]User, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteVal, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var users []User
	err = json.Unmarshal(byteVal, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// find user based on their user ID, returns user struct
func findUser(userID int) User {
	db, _ := LoadUserData()
	for _, user := range db.Users {
		if user.UserID == userID {
			return user
		}
	}

	return User{}

}

func hashResourcePath(input string) string {
	hasher.Write([]byte(input))
	return strconv.FormatUint(uint64(hasher.Sum32()), 10)
}
