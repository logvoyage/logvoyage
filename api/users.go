package main

import (
	"bitbucket.org/firstrow/logvoyage/models"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gopkg.in/kataras/iris.v6"
)

type newUser struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u newUser) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required, validation.Length(3, 255)),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 255)),
	)
}

func UsersCreate(ctx *iris.Context) {
	var data newUser
	ctx.ReadJSON(&data)

	err := data.Validate()

	if err != nil {
		response.Error(ctx, err)
		return
	}

	exists, err := models.EmailExists(data.Email)

	if err != nil {
		response.Panic(ctx, err)
		return
	}

	if exists == true {
		response.Error(ctx, "Email is already used")
		return
	}

	err = models.CreateUser(data.Email, data.Name, data.Password)

	if err != nil {
		response.Panic(ctx, err)
		return
	}

	response.Success(ctx)
}
