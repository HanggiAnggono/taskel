package main

import (
	"taskel/db"
	handler "taskel/handlers"
	model "taskel/models"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	db.DB.AutoMigrate(&model.Task{}, &model.User{})
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")

	taskHandler := handler.TaskHandler{}

	r.GET("/db/seed", func(ctx *gin.Context) {
		db.Seed()
	})
	r.GET("/db/reset", func(ctx *gin.Context) {
		db.Reset()
	})

	r.GET("/api/task/list", taskHandler.List)
	r.GET("/api/task/:id", taskHandler.Show)
	r.POST("/api/task/:id/assign", taskHandler.AssignUserToTask)
	r.POST("/api/task/:id/unassign", taskHandler.AssignUserToTask)
	// endpoint to transition task status
	r.POST("/api/task/:id/transition", taskHandler.TransitionTask)

	authHandler := handler.AuthHandler{}
	r.POST("/api/login", authHandler.Login)

	taskViewHandler := handler.TaskViewHandler{}
	r.GET("/login", authHandler.LoginView)
	r.GET("/", taskViewHandler.List)

	r.Run()
}
