package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64
	Email    string
	Name     string
	Password string
}

var (
	ErrorUserNotFound = errors.New("User not found")
)

func FindUserByEmail(email string) (*User, error) {
	var user User
	err := db.Model(&user).Where("email = ?", email).Select()
	if err != nil {
		return nil, ErrorUserNotFound
	}
	return &user, nil
}

func FindUserById(id interface{}) (*User, error) {
	var user User
	err := db.Model(&user).Where("id = ?", id).Select()
	if err != nil {
		return nil, ErrorUserNotFound
	}
	return &user, nil
}

func EmailExists(email string) (bool, error) {
	c, err := db.Model(&User{}).Where("email = ?", email).Count()
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

func CreateUser(email, name, password string) error {
	user := &User{
		Email:    email,
		Name:     name,
		Password: HashPlainPassword(password),
	}
	return db.Insert(user)
}

func HashPlainPassword(password string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b)
}
