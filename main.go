package main

import (
	"net/http"
	"taskel/config"
	"taskel/db"
	handler "taskel/handlers"
	model "taskel/models"
	"taskel/view"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db.Connect()
	db.DB.AutoMigrate(&model.Task{}, &model.User{})
	r := gin.Default()
	r.Static("/js", "./js")
	r.Use(IsAuthenticatedMiddleware())

	view.Init(r)

	taskHandler := handler.TaskHandler{}
	authHandler := handler.AuthHandler{}
	taskViewHandler := handler.TaskViewHandler{}
	// r.LoadHTMLGlob("templates/**/*")
	api := r.Group("/api")
	app := r.Group("/")
	app.Use(authorizeJWT())
	api.Use(authorizeJWT())

	r.GET("/db/seed", func(ctx *gin.Context) {
		db.Seed()
	})
	r.GET("/db/reset", func(ctx *gin.Context) {
		db.Reset()
	})

	api.GET("/task/list", taskHandler.List)
	api.GET("/task/:key", taskHandler.Show)
	api.POST("/task", taskHandler.Create)
	api.POST("/task/:key/assign", taskHandler.AssignUserToTask)
	api.POST("/task/:key/unassign", taskHandler.AssignUserToTask)
	// endpoint to transition task status
	api.POST("/task/:key/transition", taskHandler.TransitionTask)
	api.POST("/task/:key/watch", taskHandler.WatchTask)

	r.POST("/api/login", authHandler.Login)

	r.GET("/login", authHandler.LoginView)
	app.POST("/login", authHandler.Login)
	app.POST("/logout", authHandler.Logout)
	app.GET("/", taskViewHandler.List)
	app.GET("/task/new", taskViewHandler.Create)
	app.POST("/task/new", taskHandler.Create)
	app.GET("/task/:id", taskViewHandler.Show)
	app.GET("/task/:id/edit", taskViewHandler.Edit)
	app.PUT("/task/:id/edit", taskHandler.Edit)
	app.POST("/task/:id/edit", taskHandler.Edit)

	r.Run()
}

func authorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")

		abort := func() {
			contentType := c.Request.Header.Get("Content-Type")
			if contentType == "application/json" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  "error",
					"message": "unauthorized",
				})
			} else {
				c.Redirect(http.StatusMovedPermanently, "/login")
			}
			return
		}

		if err != nil || tokenStr == "" {
			abort()
			return
		}

		tokenParsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Config.JWTSecret), nil
		})

		claims, ok := tokenParsed.Claims.(jwt.MapClaims)

		if !tokenParsed.Valid || err != nil || !ok {
			abort()
			return
		}

		c.Set("jwtToken", tokenParsed)
		c.Set("userId", claims["userId"])
		view.JetView.AddGlobal("UserID", claims["userId"])
		c.Next()
	}
}

func IsAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie("token")
		view.JetView.AddGlobal("IsAuthenticated", func() bool { return token != "" })
		c.Next()
	}
}
