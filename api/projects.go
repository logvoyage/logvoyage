package main

import iris "gopkg.in/kataras/iris.v6"

func ProjectsCreate(ctx *iris.Context) {
	response.Success(ctx, map[string]string{"project": "created"})
}

func ProjectsList(ctx *iris.Context) {
	response.Success(ctx, map[string]string{"projects": "list"})
}
