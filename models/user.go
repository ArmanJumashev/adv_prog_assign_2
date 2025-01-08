package models

type User struct {
	FullName    string `json:"full_name" db:"full_name"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"`
}
