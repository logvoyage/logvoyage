package api

import (
	"github.com/logvoyage/logvoyage/elastic"
	"github.com/logvoyage/logvoyage/models"

	"gopkg.in/gin-gonic/gin.v1"
)

type projectData struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func projectsIndex(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	projects, res := models.FindAllProjectsByUser(user)

	if res.Error != nil {
		response.Panic(c, res.Error)
		return
	}

	response.Success(c, projects)
}

func projectsCreate(c *gin.Context) {
	var json projectData
	err := c.BindJSON(&json)
	if err != nil {
		response.Error(c, err)
		return
	}

	project, res := models.CreateProject(json.Name, c.MustGet("user").(*models.User))

	if res.Error != nil {
		response.Panic(c, res.Error)
		return
	}

	response.Success(c, project)
}

func projectsUpdate(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	id := c.Param("id")

	project, res := models.FindProjectById(id, user)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(c, "Project not found")
		} else {
			response.Panic(c, res.Error)
		}
		return
	}

	var json projectData
	err := c.BindJSON(json)

	if err != nil {
		response.Error(c, err)
		return
	}

	// Assign new attributes
	project.Name = json.Name

	_, res = models.SaveProject(project)

	if res.Error != nil {
		response.Panic(c, res.Error)
		return
	}

	response.Success(c, project)
}

func projectsLoad(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	id := c.Param("id")
	project, res := models.FindProjectById(id, user)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(c, "Project not found")
		} else {
			response.Panic(c, res.Error)
		}
		return
	}

	response.Success(c, project)
}

func projectsDelete(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	id := c.Param("id")
	project, res := models.FindProjectById(id, user)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(c, "Project not found")
		} else {
			response.Panic(c, res.Error)
		}
		return
	}

	res = models.DeleteProject(project)

	if res.Error != nil {
		response.Panic(c, res.Error)
		return
	}

	response.Success(c)
}

type searchQuery struct {
	Query string   `json:"query"`
	Page  int      `json:"page"`
	Types []string `json:"types"`
}

// Search log records in ElasticSearch.
func projectsLogs(c *gin.Context) {
	var query searchQuery
	err := c.BindJSON(&query)

	if err != nil {
		response.Error(c, "Error parsing request body")
		return
	}

	user := c.MustGet("user").(*models.User)
	id := c.Param("id")
	project, res := models.FindProjectById(id, user)

	if res.Error != nil {
		if res.RecordNotFound() {
			response.Error(c, "Project not found")
		} else {
			response.Panic(c, res.Error)
		}
		return
	}

	logs, err := elastic.SearchLogs(user, project, query.Types, query.Query, query.Page)

	if err != nil {
		response.Panic(c, err)
		return
	}

	response.Success(c, logs)
}

func projectsTypes(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	id := c.Param("id")
	project, _ := models.FindProjectById(id, user)
	types := elastic.GetIndexTypes(project)
	response.Success(c, types)
}
