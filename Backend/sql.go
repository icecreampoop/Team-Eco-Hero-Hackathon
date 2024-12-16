//package backend
//
//import (
//	"fmt"
//	"log"
//	"os"
//
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//)
//
//// DB is the database connection object
//var DB *gorm.DB
//
//// User struct represents the USER table schema
//type User struct {
//	UserID    int    `gorm:"primaryKey;autoIncrement"`
//	Password  string `gorm:"not null"`
//	Email     string `gorm:"unique"`
//	EXP       int    `gorm:"default:0"`
//	Level     int    `gorm:"default:1"`
//	AvatarPic string
//	Admin     bool   `gorm:"default:false"`
//	Items     []Item `gorm:"foreignKey:UserID"`
//}
//
//// Item struct represents the ITEMS table schema
//type Item struct {
//	ItemID          int    `gorm:"primaryKey;autoIncrement"`
//	UserID          int    `gorm:"not null"`
//	ItemName        string `gorm:"not null"`
//	ItemDescription string
//	Category1       string
//	Category2       string
//	Category3       string
//}
//
//// Connect function to initialize the DB connection
//func Connect() {
//	// Replace with your database credentials
//	dbUser := os.Getenv("DB_USER")
//	dbPassword := os.Getenv("DB_PASSWORD")
//	dbHost := os.Getenv("DB_HOST")
//	dbName := os.Getenv("DB_NAME")
//
//	// DSN (Data Source Name)
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName)
//
//	// Connect to MySQL
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatalf("Failed to connect to database: %v", err)
//	}
//
//	// Auto-migrate database schema
//	err = DB.AutoMigrate(&User{}, &Item{})
//	if err != nil {
//		log.Fatalf("Failed to migrate database: %v", err)
//	}
//
//	log.Println("Database connection successful!")
//}
//

package backend

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type User struct {
	UserID    int    `json:"UserID"`
	Password  string `json:"Password"`
	Email     string `json:"Email"`
	EXP       int    `json:"EXP"`
	Level     int    `json:"Level"`
	AvatarPic string `json:"AvatarPic"`
	Admin     bool   `json:"Admin"`
}

type Item struct {
	ItemID          int    `json:"ItemID"`
	UserID          int    `json:"UserID"`
	ItemName        string `json:"ItemName"`
	ItemDescription string `json:"ItemDescription"`
	Category1       string `json:"Category1"`
	Category2       string `json:"Category2"`
	Category3       string `json:"Category3"`
}

type ConfigData struct {
	Users []User `json:"users"`
	Items []Item `json:"items"`
}

var Users []User
var Items []Item

func LoadDataFromConfig(filePath string) {
	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// Parse the JSON
	var config ConfigData
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Populate global variables
	Users = config.Users
	Items = config.Items

	log.Println("Data loaded successfully from config file!")
}
