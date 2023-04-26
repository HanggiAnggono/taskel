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
	db.DB.Find(&tasks)

	c.HTML(http.StatusOK, "task-list.html", gin.H{
		"tasks": tasks,
	})
}
