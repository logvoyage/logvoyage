package main

import (
	"testing"

	"github.com/logvoyage/logvoyage/models"
	"gopkg.in/kataras/iris.v6/httptest"
)

func TestSuccessUsersCreate(t *testing.T) {
	models.GetConnection().Exec("DELETE FROM users")

	e := httptest.New(app, t)
	data := userData{
		Email:    "user@example.com",
		Name:     "test",
		Password: "password",
	}
	expected := map[string]bool{
		"success": true,
	}
	e.POST("/api/users/").WithJSON(data).Expect().JSON().Equal(expected)
}

func TestSuccessUsersLogin(t *testing.T) {
	models.GetConnection().Exec("DELETE FROM users")

	models.CreateUser("tester@example.com", "tester", "password")

	e := httptest.New(app, t)
	data := userData{
		Email:    "tester@example.com",
		Password: "password",
	}
	r := e.POST("/api/users/login").WithJSON(data).
		Expect().
		JSON().
		Object()
	r.Value("data").Object().Value("token").NotNull()
}
