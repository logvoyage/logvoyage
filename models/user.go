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

func FindUserByEmail(email string) (*User, error) {
	var user User
	res := db.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RecordNotFound() {
		return nil, ErrorUserNotFound
	}
	return &user, nil
}

func FindUserById(id interface{}) (*User, error) {
	var user User
	res := db.Where("id = ?", id).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RecordNotFound() {
		return nil, ErrorUserNotFound
	}
	return &user, nil
}

func EmailExists(email string) (bool, error) {
	count := 0
	res := db.Where("email = ?", email).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func CreateUser(email, name, password string) (*User, error) {
	user := &User{
		Email:    email,
		Name:     name,
		Password: HashPlainPassword(password),
	}
	res := db.Create(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func HashPlainPassword(password string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b)
}
