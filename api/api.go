package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/logvoyage/logvoyage/models"
	"github.com/logvoyage/logvoyage/shared/config"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/gorillamux"
)

var (
	app      *iris.Framework
	response Response
)

type Response struct {
}

// Success responses with 200 code.
// Note: only fist data argument will be passed to the response.
func (r Response) Success(ctx *iris.Context, body ...interface{}) {
	if len(body) > 0 {
		ctx.JSON(200, map[string]interface{}{"success": true, "data": body[0]})
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
	// TODO: Report error.
	log.Println("Panic:", err.Error())
	ctx.JSON(503, map[string]string{"errors": "There was an error performing your request."})
}

// Forbidden responses with 401 code, means user does not valid credentials to access handler.
func (r Response) Forbidden(ctx *iris.Context) {
	ctx.StopExecution()
	ctx.JSON(401, map[string]string{"errors": "Authentication failed"})
}

// authMiddleware performs authentication
func authMiddleware(ctx *iris.Context) {
	tokenString := ctx.RequestHeader("X-Authentication")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Get("secret")), nil
	})

	if err != nil {
		response.Forbidden(ctx)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user, res := models.FindUserByID(claims["user_id"])

		if res.Error != nil {
			if res.RecordNotFound() {
				response.Forbidden(ctx)
			} else {
				response.Error(ctx, "User not found")
			}
			return
		}

		ctx.Set("user", user)
	} else {
		response.Forbidden(ctx)
		return
	}

	ctx.Next()
}

func newCorsAdapter() iris.RouterWrapperPolicy {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "x-authentication, origin, content-type, accept, x-xsrf-token")
		w.Header().Add("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Add("Content-Type", "application/json")
		if r.Method != "OPTIONS" {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func Start(host, port string) {
	response = Response{}

	app = iris.New()

	app.Adapt(
		gorillamux.New(),
		iris.DevLogger(),
		newCorsAdapter(),
	)

	root := app.Party("/api")
	{
		userAPI := root.Party("/users")
		{
			userAPI.Post("/", UsersCreate)
			userAPI.Post("/login", UsersLogin)
		}

		projectsAPI := root.Party("/projects", authMiddleware)
		{
			projectsAPI.Get("/{id:[0-9]+}", projectsLoad)
			projectsAPI.Post("/{id:[0-9]+}", projectsUpdate)
			projectsAPI.Delete("/{id:[0-9]+}", projectsDelete)
			projectsAPI.Get("", projectsIndex)
			projectsAPI.Post("", projectsCreate)
			projectsAPI.Post("/{id:[0-9]+}/logs", projectsLogs)
			projectsAPI.Get("/{id:[0-9]+}/types", projectsTypes)
		}
	}

	dsn := fmt.Sprintf("%s:%s", host, port)
	app.Listen(dsn)
}
