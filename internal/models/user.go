package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Email string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

type UserViewModel struct {
	Id uint `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
}

func GetUser(id string) (UserViewModel, error) {
	user := User{}
	err := Database.Model(User{}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return UserViewModel{}, err
	}
	return user.mapToView(), err
}

func GetUsers() ([]UserViewModel, error) {
	users := []User{}
	viewModelUsers := []UserViewModel{}
	err := Database.Find(&users).Error
	if err != nil {
		return viewModelUsers, err
	}
	for _,user := range users {
	viewModelUsers = append(viewModelUsers, user.mapToView())
	}
	return viewModelUsers, nil
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

	return user.mapToView(), nil	
}

func (u *User) mapToView() UserViewModel{
	return UserViewModel{Id: u.ID, Username: u.Username, Email: u.Email}
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