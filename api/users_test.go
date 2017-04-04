package main

import (
	"testing"

	"bitbucket.org/firstrow/logvoyage/models"
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
	e.POST("/users/").WithJSON(data).Expect().JSON().Equal(expected)
}

func TestSuccessUsersLogin(t *testing.T) {
	models.GetConnection().Exec("DELETE FROM users")

	err := models.CreateUser("tester@example.com", "tester", "password")

	if err != nil {
		t.Error("Create user error")
	}

	e := httptest.New(app, t)
	data := userData{
		Email:    "tester@example.com",
		Password: "password",
	}
	r := e.POST("/users/login").WithJSON(data).
		Expect().
		JSON().
		Object()

	r.Value("body").Object().Value("token").NotNull()
}
