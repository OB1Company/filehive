package models

import "gorm.io/gorm"

// User contains all information for each user.
type User struct {
	gorm.Model
	Email          string `gorm:"uniqueIndex"`
	Name           string
	Salt           []byte
	HashedPassword []byte
	Country        string
	AvatarFilename string
}
