package entities

import (
	"html"
	"strings"

	"github.com/dj-ph3luy/go-playground/internal/views"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255;not null;unique" json:"username"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
	IsAdmin  bool   `gorm:"not null;default:false" json:"isAdmin"`
}

func (u *User) ToView() views.User {
	return views.User{Id: u.ID, Name: u.Name, Email: u.Email, IsAdmin: u.IsAdmin}
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))

	return nil
}
