package models

import (
	"html"
	"strings"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Email string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
	IsAdmin bool `gorm:"not null;default:false" json:"isAdmin"`
}

type UserViewModel struct {
	Id uint `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	IsAdmin bool `json:"isAdmin"`
}

func GetUser(id string) (UserViewModel, error) {
	user := User{}
	err := Database.Model(User{}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return UserViewModel{}, err
	}
	return user.MapToView(), err
}

func GetUsers() ([]UserViewModel, error) {
	users := []User{}
	viewModelUsers := []UserViewModel{}
	err := Database.Find(&users).Error
	if err != nil {
		return viewModelUsers, err
	}
	for _,user := range users {
	viewModelUsers = append(viewModelUsers, user.MapToView())
	}
	return viewModelUsers, nil
}

func (u *User) Update() error {
	err := u.beforeSave()
	if err != nil {
		return err;
	}
	return Database.Updates(&u).Error
}

func (u *User) Save() error {
	err := u.beforeSave()
	if err != nil {
		return err;
	}
	return Database.Create(&u).Error
}

func CheckLogin(username string, password string) (UserViewModel, error) {
	user := User{}
	err := Database.Model(User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return UserViewModel{}, err
	}

	err = verifyPassword(password, user.Password)
	if err != nil {
		return UserViewModel{}, err
	}

	return user.MapToView(), nil	
}

func (u *User) MapToView() UserViewModel{
	return UserViewModel{Id: u.ID, Username: u.Username, Email: u.Email, IsAdmin: u.IsAdmin}
}

func (u *User) beforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

func verifyPassword(password,hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}