package backend

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
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

	// Create a slice to hold the filtered items with OwnerUsername
	var filteredItemsWithOwner []ItemWithOwner

	for _, item := range data.Items {
		// Only include items with status "available" or "pending"
		if item.ItemStatus == "available" || item.ItemStatus == "pending" {
			// Use the findUser function to get the owner's username
			owner := findUserTpl(item.OwnerID, data.Users)

			// Create an ItemWithOwner struct and assign the OwnerUsername
			itemWithOwner := ItemWithOwner{
				Item:          item,
				OwnerUsername: owner.Username, // Assign the owner's username
			}

			// Add the item with owner information to the filtered list
			filteredItemsWithOwner = append(filteredItemsWithOwner, itemWithOwner)
		}
	}

	// Pass the filtered items with owner info to the template
	err = tpl.ExecuteTemplate(w, "items.html", filteredItemsWithOwner)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func showSingleItem(w http.ResponseWriter, r *http.Request) {
	// Extract itemID from URL parameters
	params := mux.Vars(r)
	itemID, err := strconv.Atoi(params["itemID"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	fmt.Println("Item ID:", itemID)

	// Load user data
	data, err := LoadUserData()
	if err != nil {
		fmt.Println("Error loading data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Find the item by ID
	var foundItem Item
	var itemFound bool
	for _, item := range data.Items {
		if item.ItemID == itemID {
			foundItem = item
			itemFound = true
			break
		}
	}

	// If item not found, return 404 error
	if !itemFound {
		fmt.Println("Item not found:", itemID) // Debugging line
		http.NotFound(w, r)
		return
	}

	// Render the template with the found item
	err = tpl.ExecuteTemplate(w, "item.html", foundItem)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		fmt.Println("Template execution error:", err)
	}
}

func serveUpdateItemPage(w http.ResponseWriter, r *http.Request) {

}

func createNewItemPage(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "add-item.html", nil)
	if err != nil {
		http.Error(w, "Error rendering add-item template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
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
	hashedFileName := hashResourcePath(findUser(userID).Email+r.FormValue("item-name")) + getFileExtension(handler.Header.Get("Content-Type"))
	fileResourcePath, _ := UploadFile(hashedFileName, imageBytes)
	fmt.Println(fileResourcePath)
	// add item entry to db
	userIDInt, _ := getUserID(r)
	fmt.Println(userIDInt)
	err = AddNewItem(3, r.FormValue("item-name"), r.FormValue("item-description"),
		r.FormValue("category"), fileResourcePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Respond to the client
	fmt.Println("settle")
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("File uploaded and processed successfully"))
}

func requestItem(w http.ResponseWriter, r *http.Request) {

}

func acceptRequest(w http.ResponseWriter, r *http.Request) {

}

func serveAcceptRequestPage(w http.ResponseWriter, r *http.Request) {

}

func updateItemDetails(w http.ResponseWriter, r *http.Request) {
	//asf
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	itemIDStr := r.PathValue("itemID")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	err = DeleteItem(itemID)
	if err != nil {
		if err.Error() == fmt.Sprintf("item with ID %d not found", itemID) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		log.Println("Error deleting items:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Item with ID %d successfully deleted", itemID)))

	http.Redirect(w, r, "/", http.StatusAccepted)
}

// functions to handle HTTP requests for page loads
func HandleHTTPIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Frontend/static/index.gohtml")
}

func HandleHTTPUser(w http.ResponseWriter, r *http.Request) {
	userID, exists := getUserID(r)
	if !exists {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		path := fmt.Sprintf("/user/%d", userID)
		http.Redirect(w, r, path, http.StatusSeeOther)
	}
}

func HandleHTTPSingleUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["userid"])

	// redirect function

	// userID, err := getUserID(r)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

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

func myRequestsPage(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "transactions.html", nil)
	if err != nil {
		http.Error(w, "Error rendering My-requests template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func HandleHTTPBoard(w http.ResponseWriter, r *http.Request) {
	data, err := LoadUserData()
	if err != nil {
		fmt.Println("Error loading data from JSON")
		return
	}

	users := data.Users

	sort.Slice(users, func(i, j int) bool {
		user1 := (users[i].Level * 100) + users[i].EXP
		user2 := (users[j].Level * 100) + users[j].EXP
		return user1 > user2
	})

	topFive := make(map[int]User)
	for i := 0; i < len(users) && i < 5; i++ {
		topFive[i+1] = users[i]
	}

	err = tpl.ExecuteTemplate(w, "board.html", topFive)
	if err != nil {
		http.Error(w, "Error rendering Board template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

// // HandleHTTPLogin serves the login page
// func HandleHTTPLogin(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./Frontend/static/login.html")
// }

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
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate credentials
		valid, err := ValidateUserCredentials(email, password)
		if err != nil {
			http.Error(w, "Unable to validate user credentials", http.StatusInternalServerError)
			return
		}

		if valid {
			// Set the "UserID" cookie
			userID, err := GetUserID(email)
			if err != nil {
				http.Error(w, "Unable to get user ID", http.StatusInternalServerError)
				return
			}
			cookie := http.Cookie{
				Name:  "UserID",
				Value: strconv.Itoa(userID),
			}
			http.SetCookie(w, &cookie)
			// Successful login
			http.Redirect(w, r, "/user", http.StatusSeeOther)
			return
		}

		// Invalid credentials
		tmpl.Execute(w, map[string]interface{}{
			"ErrorMessage": "Invalid login email or password",
		})
		return
	}

	// Serve the login page for GET requests
	tmpl.Execute(w, nil)
}

// HandleHTTPLogout logs the user out by deleting the "UserID" cookie and redirect to login
func HandleHTTPLogout(w http.ResponseWriter, r *http.Request) {
	// Delete the "UserID" cookie
	cookie := http.Cookie{
		Name:   "UserID",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	mux := mux.NewRouter()

	// Default handler
	mux.HandleFunc("/", showAllItems).Methods("GET") // default handler to showallitems

	mux.HandleFunc("/items", showAllItems).Methods("GET")
	mux.HandleFunc("/items/{itemID}", showSingleItem).Methods("GET")
	mux.HandleFunc("/create-item", createNewItemPage).Methods("GET")
	mux.HandleFunc("/create-item", createNewItem).Methods("POST")
	mux.HandleFunc("/items/{itemID}/request", requestItem).Methods("POST")
	mux.HandleFunc("/items/{itemID}/accept", acceptRequest).Methods("POST")
	mux.HandleFunc("/items/{itemID}/accept", serveAcceptRequestPage).Methods("GET")
	mux.HandleFunc("/items/{itemID}/update-item", updateItemDetails).Methods("POST")
	mux.HandleFunc("/items/{itemID}/update-item", serveUpdateItemPage).Methods("GET")
	mux.HandleFunc("/items/{itemID}", deleteItem).Methods("DELETE")

	mux.HandleFunc("/user", HandleHTTPUser).Methods("GET")
	mux.HandleFunc("/user/{userid}", HandleHTTPSingleUser).Methods("GET")
	mux.HandleFunc("/board", HandleHTTPBoard).Methods("GET")
	mux.HandleFunc("/login", HandleHTTPLogin).Methods("GET", "POST")
	mux.HandleFunc("/signup", HandleHTTPSignup).Methods("GET", "POST")
	mux.HandleFunc("/logout", HandleHTTPLogout).Methods("GET")
	mux.HandleFunc("/my-requests", myRequestsPage).Methods("GET")

	// Serve static files from the frontend directory
	fs := http.FileServer(http.Dir("./Frontend/static"))
	mux.PathPrefix("/Frontend/static/").Handler(http.StripPrefix("/Frontend/static/", fs))

	// Start server
	port := ":5000"
	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func getUserID(r *http.Request) (int, bool) {
	// Retrieve userID from cookie
	cookie, err := r.Cookie("UserID")
	if err != nil {
		return -1, false
	}

	userID, _ := strconv.Atoi(cookie.Value)

	return userID, true
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

// Find user based on their user ID, returns the user struct
func findUserTpl(userID int, users []User) User {
	for _, user := range users {
		if user.UserID == userID {
			return user
		}
	}
	return User{} // Return an empty user if not found
}

func hashResourcePath(input string) string {
	hasher.Write([]byte(input))
	return strconv.FormatUint(uint64(hasher.Sum32()), 10)
}

func getFileExtension(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/bmp":
		return ".bmp"
	case "image/webp":
		return ".webp"
	default:
		return "" // Unknown type
	}
}
