package backend

import (
	"fmt"
	"log"
)

// AddPoints adds points to a user's EXP and adjusts their level if necessary
func AddPoints(userID int, pointEvent int) error {
	var user User

	// Find the user by ID
	for i := range Users {
		if Users[i].UserID == userID {
			user = Users[i]
			break
		}
	}

	if user.UserID == 0 {
		log.Printf("User not found: %d", userID)
		return fmt.Errorf("user not found")
	}

	// Add points to user's EXP
	switch pointEvent {
	// Adds XP points to user
	// pointEvent 1 = 10 points (receive item)
	// pointEvent 2 = 50 points (donate item)
	case 1:
		user.EXP += 10
	case 2:
		user.EXP += 50
	}

	// Adjust user's level based on EXP
	for user.EXP >= user.Level*100 {
		user.EXP -= user.Level * 100
		user.Level++
	}

	// Update the user in the Users slice
	for i := range Users {
		if Users[i].UserID == userID {
			Users[i] = user
			break
		}
	}
	return nil
}
