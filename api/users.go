package api

import (
	"time"

	"github.com/logvoyage/logvoyage/models"
	"github.com/logvoyage/logvoyage/shared/config"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
)

type newUserData struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type loginData struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func UsersCreate(ctx *gin.Context) {
	var data newUserData
	err := ctx.BindJSON(&data)

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

func UsersLogin(ctx *gin.Context) {
	var data loginData
	err := ctx.BindJSON(&data)

	user, res := models.FindUserByEmail(data.Email)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(ctx, "User not found")
		} else {
			response.Panic(ctx, res.Error)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))

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
