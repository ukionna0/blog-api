package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password string    `gorm:"type:varchar(255);not null" json:"-"`
	Articles []Article `gorm:"foreignKey:UserID" json:"articles"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments"`
}

// HashPassword 加密密码
func (u *User) HashPassword() error {
	// 确保密码不是空的
	if u.Password == "" {
		return fmt.Errorf("password is empty")
	}

	// 检查密码长度
	if len(u.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt failed: %w", err)
	}

	u.Password = string(hashedPassword)
	fmt.Printf("DEBUG: 密码加密成功，哈希长度: %d\n", len(u.Password))
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) error {
	if u.Password == "" {
		return fmt.Errorf("stored password hash is empty")
	}
	if password == "" {
		return fmt.Errorf("provided password is empty")
	}

	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// 密码强度验证
func (u *User) ValidatePassword() error {
	if len(u.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}
