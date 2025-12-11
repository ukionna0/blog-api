package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title    string    `gorm:"type:varchar(255);not null" json:"title"`
	Content  string    `gorm:"type:text;not null" json:"content"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	Comments []Comment `gorm:"foreignKey:ArticleID" json:"comments"`
}
