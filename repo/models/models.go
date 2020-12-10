package models

import "gorm.io/gorm"

// User contains all information for each user.
type User struct {
	gorm.Model
	ID              string `json:"id" gorm:"primary_key"`
	Email           string `gorm:"uniqueIndex"`
	Name            string
	Salt            []byte
	HashedPassword  []byte
	Country         string
	AvatarFilename  string
	FilecoinAddress string
}

// Dataset holds metadata about a dataaset.
type Dataset struct {
	gorm.Model
	ID               string  `json:"id" gorm:"primary_key"`
	UserID           string  `json:"userID"`
	Title            string  `json:"title"`
	ShortDescription string  `json:"shortDescription"`
	FullDescription  string  `json:"fullDescription"`
	ImageFilename    string  `json:"imageFilename"`
	FileType         string  `json:"fileType"`
	Price            float64 `json:"price"`
}
