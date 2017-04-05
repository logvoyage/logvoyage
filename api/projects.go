package main

import (
	"bitbucket.org/firstrow/logvoyage/models"

	validation "github.com/go-ozzo/ozzo-validation"
	iris "gopkg.in/kataras/iris.v6"
)

type projectData struct {
	Name string `json:"name"`
}

func (p projectData) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(1, 255)),
	)

}

func ProjectsCreate(ctx *iris.Context) {
	var data projectData
	ctx.ReadJSON(&data)

	err := data.Validate()

	if err != nil {
		response.Error(ctx, err)
		return
	}

	project, err := models.CreateProject(data.Name, ctx.Get("user").(*models.User))

	if err != nil {
		response.Panic(ctx, err)
		return
	}

	response.Success(ctx, project)
}

func ProjectsList(ctx *iris.Context) {
	user := ctx.Get("user").(*models.User)
	projects, err := models.FindAllProjectsByUser(user)

	if err != nil {
		response.Panic(ctx, err)
		return
	}

	response.Success(ctx, map[string]interface{}{"projects": projects})
}
