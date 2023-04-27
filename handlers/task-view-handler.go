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
	db.DB.Find(&tasks, "user_id = ?", c.MustGet("userId"))

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

func (h *TaskViewHandler) Create(c *gin.Context) {
	var users []model.User
	db.DB.Find(&users)

	c.HTML(http.StatusOK, "tasks/create", gin.H{
		"title":      "Create Task",
		"taskStatus": []string{"todo", "inprogress", "done"},
		"users":      users,
	})
}
