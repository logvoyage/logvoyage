package main

import (
	"time"

	"bitbucket.org/firstrow/logvoyage/models"
	"bitbucket.org/firstrow/logvoyage/shared/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/kataras/iris.v6"
)

type userData struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u userData) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required, validation.Length(3, 255)),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 255)),
	)
}

func UsersCreate(ctx *iris.Context) {
	var data userData
	ctx.ReadJSON(&data)

	err := data.Validate()

	if err != nil {
		response.Error(ctx, err)
		return
	}

	exists, res := models.EmailExists(data.Email)

	if res.Error != nil {
		response.Panic(ctx, res.Error)
		return
	}

	if exists {
		response.Error(ctx, "Email is already used")
		return
	}

	_, res = models.CreateUser(data.Email, data.Name, data.Password)

	if res.Error != nil {
		response.Panic(ctx, res.Error)
		return
	}

	response.Success(ctx)
}

func UsersLogin(ctx *iris.Context) {
	var data userData
	ctx.ReadJSON(&data)

	user, res := models.FindUserByEmail(data.Email)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(ctx, "User not found")
		} else {
			response.Panic(ctx, res.Error)
		}
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))

	if err != nil {
		response.Error(ctx, "User not found")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"timestamp": time.Now().UTC().Unix(),
	})

	secret := []byte(config.Get("secret"))
	tokenString, err := token.SignedString(secret)

	if err != nil {
		response.Panic(ctx, err)
		return
	}

	response.Success(ctx, map[string]string{"token": tokenString})
}
