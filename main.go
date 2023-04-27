package main

import (
	"fmt"
	"net/http"
	"taskel/config"
	"taskel/db"
	handler "taskel/handlers"
	model "taskel/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db.Connect()
	db.DB.AutoMigrate(&model.Task{}, &model.User{})
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	api := r.Group("/api")
	api.Use(authorizeJWT())

	taskHandler := handler.TaskHandler{}

	r.GET("/db/seed", func(ctx *gin.Context) {
		db.Seed()
	})
	r.GET("/db/reset", func(ctx *gin.Context) {
		db.Reset()
	})

	api.GET("/task/list", taskHandler.List)
	api.GET("/task/:id", taskHandler.Show)
	api.POST("/task/:id/assign", taskHandler.AssignUserToTask)
	api.POST("/task/:id/unassign", taskHandler.AssignUserToTask)
	// endpoint to transition task status
	api.POST("/task/:id/transition", taskHandler.TransitionTask)

	authHandler := handler.AuthHandler{}
	r.POST("/api/login", authHandler.Login)

	taskViewHandler := handler.TaskViewHandler{}
	r.GET("/login", authHandler.LoginView)
	r.POST("/login", authHandler.Login)
	r.GET("/", authorizeJWT(), taskViewHandler.List)

	r.Run()
}

func authorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")

		abort := func() {
			fmt.Print("TOKEN TOKEN TOKEN TOKEN TOKEN TOKEN TOKEN TOKEN TOKEN TOKEN TOKEN TOKEN TOKEN \n", err)
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

		if !tokenParsed.Valid || err != nil {
			abort()
			return
		}

		c.Next()
	}
}
