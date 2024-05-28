package utility

import (
	"regexp"
	"ungraded-challenge-6/entity"
)

func ValidateUser(user entity.User) (bool, string) {
	if user.Email == "" || !isValidEmail(user.Email) {
		return false, "Invalid or empty email"
	}
	if len(user.Password) < 8 {
		return false, "Password must be at least 8 characters long and cannot be empty"
	}
	if len(user.FullName) < 6 || len(user.FullName) > 15 {
		return false, "Full name must be between 6 and 15 characters long and cannot be empty"
	}
	if user.Age < 17 {
		return false, "Age must be at least 17 and cannot be empty"
	}
	if user.Occupation == "" {
		return false, "Occupation cannot be empty"
	}
	if !(user.Role == "admin" || user.Role == "superadmin") {
		return false, "Role must be 'admin' or 'superadmin'"
	}
	return true, ""
}

func isValidEmail(email string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return reg.MatchString(email)
}
