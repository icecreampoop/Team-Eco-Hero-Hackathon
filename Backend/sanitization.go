package backend

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeInput sanitizes input by removing HTML and trimming spaces
func SanitizeInput(input string) string {
	re := regexp.MustCompile("<.*?>")
	input = re.ReplaceAllString(input, "")
	input = html.EscapeString(input)
	input = strings.TrimSpace(input)
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")
	return input
}

// IsValidUsername checks if the username is valid
func IsValidUsername(username string) bool {
	if len(username) < 15 || len(username) > 30 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9@.]+$", username)

	return matched
}

// IsValidPassword checks if the password is valid (at least 6 chars and contains a number)
func IsValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	matched, _ := regexp.MatchString("[a-zA-Z0-9]+$", password)
	return matched
}
