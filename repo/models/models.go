package models

// User contains all information for each user.
type User struct {
	Email          string `gorm:"primary_key"`
	Name           string
	Salt           []byte
	HashedPassword []byte
	Country        string
}
