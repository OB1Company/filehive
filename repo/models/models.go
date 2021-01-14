package models

import (
	"gorm.io/gorm"
	"time"
)

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
	gorm.Model       `json:"-"`
	ID               string    `json:"id" gorm:"primary_key"`
	CreatedAt        time.Time `gorm:"index" json:"-"`
	UserID           string    `json:"userID"`
	JobID            string    `json:"jobID"`
	Username         string    `json:"username"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"shortDescription"`
	FullDescription  string    `json:"fullDescription"`
	ImageFilename    string    `json:"imageFilename"`
	FileType         string    `json:"fileType"`
	Price            float64   `json:"price"`
	Views            int64     `json:"totalViews"`
	Purchases        int64     `json:"totalPurchases"`
}

// Purchase holds information about a user purchase.
type Purchase struct {
	gorm.Model       `json:"-"`
	ID               string `json:"id" gorm:"primary_key"`
	UserID           string `json:"userID"`
	DatasetID        string `json:"datasetID"`
	Timestamp        time.Time
	Title            string `json:"title"`
	ShortDescription string `json:"shortDescription"`
	ImageFilename    string `json:"imageFilename"`
	FileType         string `json:"fileType"`
	Username         string `json:"username"`
}

// Click represents a view on a dataset.
type Click struct {
	gorm.Model
	DatasetID string
	Timestamp time.Time `gorm:"index"`
}
