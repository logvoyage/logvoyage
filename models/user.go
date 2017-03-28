package models

import (
	"errors"
)

type User struct {
	Id       int
	Email    string
	Name     string
	Password string
}

var (
	ErrorUserNotFound = errors.New("User not found")
)

func FindUserByEmail(email string) (*User, error) {
	var user User
	err := db.Model(&user).Where("email = ?", email).Select(&user)
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
		Password: password,
	}
	return db.Insert(user)
}