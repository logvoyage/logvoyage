package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/logvoyage/logvoyage/models"
	"github.com/logvoyage/logvoyage/shared/config"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"
)

var (
	app      *gin.Engine
	response Response
)

type Response struct {
}

// Success responses with 200 code.
// Note: only first data argument will be passed to the response.
func (r Response) Success(ctx *gin.Context, body ...interface{}) {
	if len(body) > 0 {
		ctx.JSON(200, gin.H{"success": true, "data": body[0]})
	} else {
		ctx.JSON(200, gin.H{"success": true})
	}
}

// Error returns 200 OK response with json field "errors" with error descrioption.
// This function should be used to display validation or other expected errors.
// errors may be string or array of hashes.
func (r Response) Error(ctx *gin.Context, err interface{}) {
	ctx.JSON(400, gin.H{"errors": err})

}

// Panic responses with 503 error.
func (r Response) Panic(ctx *gin.Context, err error) {
	// TODO: Report error.
	log.Println("Panic:", err.Error())
	ctx.JSON(503, gin.H{"errors": "There was an error performing your request."})
}

// Forbidden responses with 401 code, means user does not valid credentials to access handler.
func (r Response) Forbidden(ctx *gin.Context) {
	ctx.Abort()
	ctx.JSON(401, gin.H{"errors": "Authentication failed"})
}

// authMiddleware performs user authentication using jwt tokens.
func authMiddleware(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("X-Authentication")

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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "x-authentication, origin, content-type, accept, x-xsrf-token")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")

		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

func InitRoutes() {
	response = Response{}

	app = gin.Default()
	nt
	app.Use(corsMiddleware())

	root := app.Group("/api")
	{
		userAPI := root.Group("/users")
		{
			userAPI.POST("", UsersCreate)
			userAPI.POST("/login", UsersLogin)
		}

		projectsAPI := root.Group("/projects", authMiddleware)
		{
			projectsAPI.GET("/:id", projectsLoad)
			projectsAPI.POST("/:id", projectsUpdate)
			projectsAPI.DELETE("/:id", projectsDelete)
			projectsAPI.GET("", projectsIndex)
			projectsAPI.POST("", projectsCreate)
			projectsAPI.POST("/:id/logs", projectsLogs)
			projectsAPI.GET("/:id/types", projectsTypes)
		}
	}
}

func Start(host, port string) {
	dsn := fmt.Sprintf("%s:%s", host, port)
	app.Run(dsn)
}
