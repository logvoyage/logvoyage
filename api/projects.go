package main

import (
	"github.com/logvoyage/logvoyage/elastic"
	"github.com/logvoyage/logvoyage/models"

	validation "github.com/go-ozzo/ozzo-validation"
	iris "gopkg.in/kataras/iris.v6"
)

type projectData struct {
	Name string
}

func (p projectData) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(1, 30)),
	)
}

func projectsIndex(ctx *iris.Context) {
	user := ctx.Get("user").(*models.User)
	projects, res := models.FindAllProjectsByUser(user)

	if res.Error != nil {
		response.Panic(ctx, res.Error)
		return
	}

	response.Success(ctx, projects)
}

func projectsCreate(ctx *iris.Context) {
	var data projectData
	ctx.ReadJSON(&data)

	err := data.Validate()

	if err != nil {
		response.Error(ctx, err)
		return
	}

	project, res := models.CreateProject(data.Name, ctx.Get("user").(*models.User))

	if res.Error != nil {
		response.Panic(ctx, res.Error)
		return
	}

	response.Success(ctx, project)
}

func projectsUpdate(ctx *iris.Context) {
	user := ctx.Get("user").(*models.User)
	id, _ := ctx.ParamInt("id")

	project, res := models.FindProjectById(id, user)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(ctx, "Project not found")
		} else {
			response.Panic(ctx, res.Error)
		}
		return
	}

	var data projectData
	ctx.ReadJSON(&data)

	err := data.Validate()

	if err != nil {
		response.Error(ctx, err)
		return
	}

	// Assign new attributes
	project.Name = data.Name

	_, res = models.SaveProject(project)

	if res.Error != nil {
		response.Panic(ctx, res.Error)
		return
	}

	response.Success(ctx, project)
}

func projectsLoad(ctx *iris.Context) {
	user := ctx.Get("user").(*models.User)
	id, _ := ctx.ParamInt("id")
	project, res := models.FindProjectById(id, user)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(ctx, "Project not found")
		} else {
			response.Panic(ctx, res.Error)
		}
		return
	}

	response.Success(ctx, project)
}

func projectsDelete(ctx *iris.Context) {
	user := ctx.Get("user").(*models.User)
	id, _ := ctx.ParamInt("id")
	project, res := models.FindProjectById(id, user)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(ctx, "Project not found")
		} else {
			response.Panic(ctx, res.Error)
		}
		return
	}

	res = models.DeleteProject(project)

	if res.Error != nil {
		response.Panic(ctx, res.Error)
		return
	}

	response.Success(ctx)
}

type searchQuery struct {
	Query string   `json:"query"`
	Page  int      `json:"page"`
	Types []string `json:"types"`
}

// Search log records in ElasticSearch.
func projectsLogs(ctx *iris.Context) {
	var query searchQuery
	err := ctx.ReadJSON(&query)

	if err != nil {
		response.Error(ctx, "Error parsing request body")
		return
	}

	user := ctx.Get("user").(*models.User)
	id, _ := ctx.ParamInt("id")
	project, res := models.FindProjectById(id, user)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(ctx, "Project not found")
		} else {
			response.Panic(ctx, res.Error)
		}
		return
	}

	logs, err := elastic.SearchLogs(user, project, query.Types, query.Query, query.Page)

	if err != nil {
		response.Panic(ctx, err)
		return
	}

	response.Success(ctx, logs)
}

func projectsTypes(ctx *iris.Context) {
	user := ctx.Get("user").(*models.User)
	id, _ := ctx.ParamInt("id")
	project, _ := models.FindProjectById(id, user)
	types := elastic.GetIndexTypes(project)
	response.Success(ctx, types)
}
