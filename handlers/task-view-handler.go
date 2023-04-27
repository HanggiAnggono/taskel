package handler

import (
	"fmt"
	"net/http"
	"taskel/db"
	model "taskel/models"

	"github.com/gin-gonic/gin"
)

type TaskViewHandler struct{}

func (h *TaskViewHandler) List(c *gin.Context) {
	var tasks []model.Task
	db.DB.Find(&tasks)

	cookies := c.Request.Cookies()
	fmt.Println("COOKIES")
	for _, cookie := range cookies {
		fmt.Println(cookie.Name, cookie.Value)
	}

	c.HTML(http.StatusOK, "tasks/index", gin.H{
		"title": "My Tasks",
		"tasks": tasks,
	})
}
