package main

import (
	"testing"

	"bitbucket.org/firstrow/logvoyage/models"
	"gopkg.in/kataras/iris.v6/httptest"
)

func TestSuccessCreateUser(t *testing.T) {
	models.GetConnection().Exec("DELETE FROM users")

	e := httptest.New(app, t)
	data := userData{
		Email:    "user@example.com",
		Name:     "test",
		Password: "123456",
	}
	expected := map[string]bool{
		"success": true,
	}
	e.POST("/users/").WithJSON(data).Expect().JSON().Equal(expected)
}
