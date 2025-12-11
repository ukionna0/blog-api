package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content   string  `gorm:"type:text;not null" json:"content"`
	UserID    uint    `gorm:"not null" json:"user_id"`
	ArticleID uint    `gorm:"not null" json:"article_id"`
	User      User    `gorm:"foreignKey:UserID" json:"user"`
	Article   Article `gorm:"foreignKey:ArticleID" json:"article"`
}
