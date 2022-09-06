package db

import (
	"gorm.io/gorm"
)

// Database Models

// User struct
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" form:"username" binding:"required" json:"username"`
	Password string `form:"password" binding:"required" json:"password"`
}

// Command struct
type Command struct {
	gorm.Model
	Keyword     string `gorm:"uniqueIndex" form:"keyword" binding:"required" json:"keyword"`
	Description string `form:"description" binding:"required" json:"description"`
	Text        string `form:"text" binding:"required" json:"text"`
	Link        string `gorm:"default:none" form:"link" json:"link"`
}
