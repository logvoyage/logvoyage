package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email    string
	Name     string
	Password string
}

var (
	ErrorUserNotFound = errors.New("User not found")
)

func FindUserByEmail(email string) (*User, *gorm.DB) {
	var user User
	res := db.Model(&user).Where("email = ?", email).First(&user)
	return &user, res
}

func FindUserByID(id interface{}) (*User, *gorm.DB) {
	var user User
	res := db.Model(&user).Where("id = ?", id).First(&user)
	return &user, res
}

func EmailExists(email string) (bool, *gorm.DB) {
	count := 0
	res := db.Model(&User{}).Where("email = ?", email).Count(&count)
	return count > 0, res
}

func CreateUser(email, name, password string) (*User, *gorm.DB) {
	user := &User{
		Email:    email,
		Name:     name,
		Password: HashPlainPassword(password),
	}
	res := db.Create(user)
	return user, res
}

func HashPlainPassword(password string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b)
}
