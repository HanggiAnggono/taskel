package handler

import (
	"net/http"
	"taskel/db"
	model "taskel/models"

	"github.com/gin-gonic/gin"
)

type TaskViewHandler struct{}

func (h *TaskViewHandler) List(c *gin.Context) {
	var tasks []model.Task
	db.DB.Preload("User").Find(&tasks)

	c.HTML(http.StatusOK, "tasks/index", gin.H{
		"title": "My Tasks",
		"tasks": tasks,
	})
}

func (h *TaskViewHandler) StatusColor(status string) string {
	switch status {
	case "inprogress":
		return "bg-warning"
	case "done":
		return "bg-success"
	default:
		return "bg-secondary"
	}
}

func (h *TaskViewHandler) Show(c *gin.Context) {
	id := c.Param("id")

	var task model.Task
	db.DB.Preload("User").First(&task, id)

	accept := c.Request.Header.Get("Accept")

	if task.ID == 0 {
		switch accept {
		case "application/json":
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "task not found",
			})
		default:
			c.AbortWithStatus(http.StatusNotFound)
		}
		return
	}

	switch accept {
	case "application/json":
		c.JSON(http.StatusOK, gin.H{
			"data":    task,
			"status":  http.StatusOK,
			"message": "success",
		})
	default:
		c.HTML(http.StatusOK, "tasks/show", gin.H{
			"title": task.Title,
			"task":  task,
		})
	}
}

func (h *TaskViewHandler) Create(c *gin.Context) {
	var users []model.User
	db.DB.Find(&users)

	c.HTML(http.StatusOK, "tasks/create", gin.H{
		"title":      "Create Task",
		"taskStatus": []string{"todo", "inprogress", "done"},
		"users":      users,
	})
}
