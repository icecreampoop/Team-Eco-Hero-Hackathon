package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type User struct {
	UserID         int    `json:"UserID"`
	Password       string `json:"Password"`
	Email          string `json:"Email"`
	EXP            int    `json:"EXP"`
	Level          int    `json:"Level"`
	AvatarPic      string `json:"AvatarPic"`
	Admin          bool   `json:"Admin"`
	ActiveRequests []int  `json:"ActiveRequests"`
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
	ItemID            int        `json:"ItemID"`
	OwnerID           int        `json:"OwnerID"`
	ReceiverID        int        `json:"ReceiverID"`
	ItemName          string     `json:"ItemName"`
	ItemDescription   string     `json:"ItemDescription"`
	Category          Categories `json:"Category"`
	ItemStatus        Statuses   `json:"Status"`
	CurrentRequesters []int      `json:"CurrentRequesters"`
}

// Data struct to represent the data.json structure
type Data struct {
	Users []User `json:"users"`
	Items []Item `json:"items"`
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
func AddNewUser(email, password string) error {
	data, err := LoadUserData()
	if err != nil {
		return err
	}

	newUser := User{
		UserID:    len(data.Users) + 1,
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
func AddNewItem(ownerID int, itemName string, categories string) error {
	data, err := LoadUserData()
	if err != nil {
		return err
	}
	userExists := false
	for _, user := range data.Users {
		if user.UserID == ownerID {
			userExists = true
			break
		}
	}
	if !userExists {
		return err
	}

	var validCategory Categories
	switch categories {
	case string(Electronics), string(Mobile), string(Furniture), string(HardwareTools),
		string(Sports), string(Clothing), string(Books), string(Media), string(Others):
		validCategory = Categories(categories)
	default:
		return fmt.Errorf("invalid category: %s", categories)
	}
	// Create a new item
	newItem := Item{
		ItemID:            len(data.Items) + 1,
		OwnerID:           ownerID,
		ReceiverID:        0, // 0 means no receiver yet
		ItemName:          itemName,
		ItemDescription:   "", // Optional: allow it to be empty for now
		Category:          validCategory,
		ItemStatus:        StatusAvailable, // New items are 'available' by default
		CurrentRequesters: []int{},
	}

	// Append the new item to the list
	data.Items = append(data.Items, newItem)

	// Save the updated data
	return SaveUserData(data)
}
