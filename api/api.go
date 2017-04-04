package main

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

var (
	app      *iris.Framework
	response Response
)

type Response struct {
}

// Success responses with 200 code.
// Note: only fist body argument will be passed to the response.
func (r Response) Success(ctx *iris.Context, body ...interface{}) {
	if len(body) > 0 {
		ctx.JSON(200, map[string]interface{}{"success": true, "body": body[0]})
	} else {
		ctx.JSON(200, map[string]interface{}{"success": true})
	}
}

// Error returns 200 OK response with json field "errors" with error descrioption.
// This function should be used to display validation or other expected errors.
// errors may be string or array of hashes.
func (r Response) Error(ctx *iris.Context, err interface{}) {
	ctx.JSON(200, map[string]interface{}{"errors": err})

}

// Panic responses with 503 error.
func (r Response) Panic(ctx *iris.Context, err error) {
	// TODO: Send orignal error to issue tracker.
	ctx.JSON(503, map[string]interface{}{"errors": err.Error()})
}

func init() {
	response = Response{}

	app = iris.New()
	app.Adapt(httprouter.New())
	app.Adapt(iris.DevLogger())

	userAPI := app.Party("/users")
	{
		userAPI.Post("/", UsersCreate)
		userAPI.Post("/login", UsersLogin)
	}
}

func main() {
	app.Listen("127.0.0.1:3000")
}
