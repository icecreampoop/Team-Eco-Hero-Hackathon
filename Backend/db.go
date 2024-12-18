package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var Users []User
var Items []Item

type User struct {
	UserID    int    `json:"UserID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	Email     string `json:"Email"`
	EXP       int    `json:"EXP"`
	Level     int    `json:"Level"`
	AvatarPic string `json:"AvatarPic"`
	Admin     bool   `json:"Admin"`
	// array of item IDs
	ActiveRequests []int `json:"ActiveRequests"`
}

type Categories string
type Statuses string

const (
	StatusAvailable Statuses   = "available"
	StatusPending   Statuses   = "pending"
	StatusDonated   Statuses   = "donated"
	Electronics     Categories = "electronics"
	Mobile          Categories = "mobile"
	Furniture       Categories = "furniture"
	HardwareTools   Categories = "hardwaretools"
	Sports          Categories = "sports"
	Clothing        Categories = "clothing"
	Books           Categories = "books"
	Media           Categories = "media"
	Others          Categories = "others"
)

type Item struct {
	ItemID          int        `json:"ItemID"`
	OwnerID         int        `json:"OwnerID"`
	ReceiverID      int        `json:"ReceiverID"`
	ItemName        string     `json:"ItemName"`
	ItemDescription string     `json:"ItemDescription"`
	ItemImageLink   string     `json:"ItemImageLink"`
	Category        Categories `json:"Category"`
	ItemStatus      Statuses   `json:"Status"`
	// array of user IDs
	CurrentRequesters []int `json:"CurrentRequesters"`
}

// Data struct to represent the data.json structure
type Data struct {
	Users []User `json:"users"`
	Items []Item `json:"items"`
}

type ItemWithOwner struct {
	Item
	OwnerUsername            string // This will hold the username of the item owner
	OwnerID                  int
	CurrentRequestersNameArr []string
}

// Load user data from data.json
func LoadUserData() (Data, error) {
	var data Data
	file, err := os.Open("data.json")
	if err != nil {
		return data, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(bytes, &data)
	return data, err
}

// Save user data to data.json
func SaveUserData(data Data) error {
	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}

// ValidateUserCredentials validates the user credentials
func ValidateUserCredentials(email, password string) (bool, error) {
	data, err := LoadUserData()
	if err != nil {
		return false, err
	}

	for _, user := range data.Users {
		if user.Email == email && user.Password == password {
			return true, nil
		}
	}

	return false, nil
}

// AddNewUser adds a new user to the data.json file
func AddNewUser(email, password, username string) error {
	data, err := LoadUserData()
	if err != nil {
		return err
	}

	newUser := User{
		UserID:    len(data.Users) + 1,
		Username:  username,
		Password:  password,
		Email:     email,
		EXP:       0,
		Level:     1,
		AvatarPic: "",
		Admin:     false,
	}

	data.Users = append(data.Users, newUser)
	return SaveUserData(data)
}

// Add new item to data.json
func AddNewItem(ownerID int, itemName string, itemDescription string, categories string, imageLink string) error {
	data, err := LoadUserData()
	if err != nil {
		return err
	}
	userExists := false
	for _, user := range data.Users {
		if user.UserID == ownerID {
			userExists = true
			fmt.Printf("user %v exists", user.UserID)
			break
		}
	}
	if !userExists {
		fmt.Println("failed no such user")
		return err
	}

	var validCategory Categories
	switch categories {
	case string(Electronics), string(Mobile), string(Furniture), string(HardwareTools),
		string(Sports), string(Clothing), string(Books), string(Media), string(Others):
		validCategory = Categories(categories)
	default:
		fmt.Println("failed cat")

		return fmt.Errorf("invalid category: %s", categories)
	}
	// Create a new item
	newItem := Item{
		ItemID:            len(data.Items) + 1,
		OwnerID:           ownerID,
		ReceiverID:        0, // 0 means no receiver yet
		ItemName:          itemName,
		ItemDescription:   itemDescription,
		Category:          validCategory,
		ItemStatus:        StatusAvailable, // New items are 'available' by default
		CurrentRequesters: []int{},
		ItemImageLink:     imageLink,
	}

	// Append the new item to the list
	data.Items = append(data.Items, newItem)
	// Save the updated data
	return SaveUserData(data)
}

// GetUserID returns the UserID of a user by their email
func GetUserID(email string) (int, error) {
	data, err := LoadUserData()
	if err != nil {
		return 0, err
	}

	for _, user := range data.Users {
		if user.Email == email {
			return user.UserID, nil
		}
	}

	return 0, fmt.Errorf("user with email %s not found", email)
}

// GetItem retrieves an item by its ItemID
func GetItem(itemID int) (Item, error) {
	var data Data
	file, err := os.Open("data.json")
	if err != nil {
		return Item{}, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return Item{}, err
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return Item{}, err
	}

	for _, item := range data.Items {
		if item.ItemID == itemID {
			return item, nil
		}
	}

	return Item{}, fmt.Errorf("item with ID %d not found", itemID)
}

// RequestItem updates the Requesters field of an item in the data.json file by its ItemID
func RequestItem(itemID int, userID int) error {
	// Load existing data
	data, err := LoadUserData()
	if err != nil {
		return err
	}

	// Iterate through the items to find the item by itemID
	for i, item := range data.Items {
		if item.ItemID == itemID {
			if item.OwnerID == userID {
				return fmt.Errorf("cannot request your own item")
			}
			// Add the userID to the requesters
			// Check for duplicates before adding the userID to CurrentRequesters
			userExists := false
			for _, requesterID := range data.Items[i].CurrentRequesters {
				if requesterID == userID {
					userExists = true
					break
				}
			}
			if !userExists {
				data.Items[i].CurrentRequesters = append(data.Items[i].CurrentRequesters, userID)
			}

			// Add itemID to ActiveRequests if not already present
			for j, user := range data.Users {
				if user.UserID == userID {
					itemExists := false
					for _, activeItemID := range data.Users[j].ActiveRequests {
						if activeItemID == itemID {
							itemExists = true
							break
						}
					}
					if !itemExists {
						data.Users[j].ActiveRequests = append(data.Users[j].ActiveRequests, itemID)
					} else {
						return fmt.Errorf("you already requested for this")
					}
				}
			}
			break
		}
	}

	// Save the updated data
	err = SaveUserData(data)
	if err != nil {
		return err
	}

	return nil
}

// DeleteItem removes an item from the data.json file by its ItemID
func DeleteItem(itemID int, requesterID int) error {
	// Load existing data
	data, err := LoadUserData()
	if err != nil {
		return err
	}

	// Flag to track if the item was found
	itemFound := false

	// Create a slice to hold the remaining items
	var updatedItems []Item

	// Iterate through the items
	for _, item := range data.Items {
		if item.ItemID == itemID {
			itemFound = true
			if requesterID != item.OwnerID {
				return fmt.Errorf("cannot delete item. your id %d, original id %d", requesterID, item.OwnerID)
			}
			continue // Skip the item to delete it
		}
		updatedItems = append(updatedItems, item)
	}

	// Check if the item was found
	if !itemFound {
		return fmt.Errorf("item with ID %d not found", itemID)
	}

	// Update the data with the remaining items
	data.Items = updatedItems

	// Save the updated data
	err = SaveUserData(data)
	if err != nil {
		return err
	}

	return nil
}

// EditItem updates the fields of an existing item in the database
func EditItem(itemID int, newName, newDescription, newLink, newCategory string) error {
	// Load existing data
	data, err := LoadUserData()
	if err != nil {
		return err
	}

	// Validate category
	if newCategory != "" {
		validCategory := false
		for _, c := range []Categories{
			Electronics, Mobile, Furniture, HardwareTools,
			Sports, Clothing, Books, Media, Others,
		} {
			if string(c) == newCategory {
				validCategory = true
				break
			}
		}
		if !validCategory {
			return fmt.Errorf("invalid category: %s", newCategory)
		}
	}

	// Flag to check if the item exists
	itemFound := false

	// Search for the item and update its fields
	for i, item := range data.Items {
		if item.ItemID == itemID {
			if newName != "" {
				data.Items[i].ItemName = newName
			}
			if newDescription != "" {
				data.Items[i].ItemDescription = newDescription
			}
			if newLink != "" {
				data.Items[i].ItemImageLink = newLink
			}
			if newCategory != "" {
				data.Items[i].Category = Categories(newCategory)
			}
			itemFound = true
			break
		}
	}

	// If the item was not found, return an error
	if !itemFound {
		return fmt.Errorf("item with ID %d not found", itemID)
	}

	// Save the updated data back to the file
	err = SaveUserData(data)
	if err != nil {
		return err
	}

	return nil
}
