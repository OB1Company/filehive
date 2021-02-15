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
	PowergateToken  string
	PowergateID     string
	ActivationCode  string
	Activated       bool `gorm:"default:false;not null"`
	ResetToken      string
	ResetValid      time.Time
	Admin           bool `gorm:"default:false;not null"`
	Disabled        bool `gorm:"default:false;not null"`
}

// Dataset holds metadata about a dataaset.
type Dataset struct {
	gorm.Model       `json:"-"`
	ID               string    `json:"id" gorm:"primary_key"`
	CreatedAt        time.Time `gorm:"index" json:"createdAt"`
	UserID           string    `json:"userID"`
	JobID            string    `json:"jobID"`
	ContentID        string    `json:"contentID"`
	Username         string    `json:"username"`
	Title            string    `gorm:"index:idx_search" json:"title"`
	ShortDescription string    `gorm:"index:idx_search" json:"shortDescription"`
	FullDescription  string    `gorm:"index:idx_search" json:"fullDescription"`
	ImageFilename    string    `json:"imageFilename"`
	DatasetFilename  string    `json:"datasetFilename"`
	FileType         string    `json:"fileType"`
	FileSize         int64     `json:"fileSize"`
	Price            float64   `json:"price"`
	Views            int64     `json:"totalViews"`
	Purchases        int64     `json:"totalPurchases"`
	Delisted         bool      `gorm:"default:false;non null"`
}

// Purchase holds information about a user purchase.
type Purchase struct {
	gorm.Model       `json:"-"`
	ID               string `json:"id" gorm:"primary_key"`
	UserID           string `json:"userID"`
	SellerID         string `json:"sellerID"`
	DatasetID        string `json:"datasetID"`
	Timestamp        time.Time
	Title            string  `json:"title"`
	ShortDescription string  `json:"shortDescription"`
	ImageFilename    string  `json:"imageFilename"`
	FileType         string  `json:"fileType"`
	Username         string  `json:"username"`
	Price            float64 `json:"price"`
	Cid              string  `json:"cid"`
}

// Click represents a view on a dataset.
type Click struct {
	gorm.Model
	DatasetID string
	Timestamp time.Time `gorm:"index"`
}
