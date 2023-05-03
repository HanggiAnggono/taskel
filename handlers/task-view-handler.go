package handler

import (
	"math"
	"net/http"
	"taskel/db"
	model "taskel/models"
	"taskel/view"

	"github.com/gin-gonic/gin"
)

type TaskViewHandler struct{}

type TaskListRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=10"`
}

func (h *TaskViewHandler) List(c *gin.Context) {
	var tasks []model.Task
	var request TaskListRequest
	c.ShouldBind(&request)
	offset := (request.Page - 1) * request.PageSize
	db.DB.Preload("User").Limit(request.PageSize).Offset(offset).Find(&tasks)

	var totalItems int64
	db.DB.Model(&model.Task{}).Count(&totalItems)
	totalPageFloat := math.Ceil(float64(totalItems) / float64(request.PageSize))
	totalPage := int(totalPageFloat)

	view.HTML(c, http.StatusOK, "tasks/index.jet", gin.H{
		"title":      "My Tasks",
		"tasks":      tasks,
		"page":       request.Page,
		"pageSize":   request.PageSize,
		"totalPage":  totalPage,
		"totalItems": totalItems,
	}, nil)
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

	var task *model.Task
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
		view.HTML(c, http.StatusOK, "tasks/show.jet", gin.H{
			"title": task.Title,
			"task":  task,
		}, nil)
	}
}

func (h *TaskViewHandler) Create(c *gin.Context) {
	var users []model.User
	db.DB.Find(&users)

	task := model.Task{}

	view.HTML(c, http.StatusOK, "tasks/create.jet", gin.H{
		"title":      "Create Task",
		"taskStatus": []string{"todo", "inprogress", "done"},
		"users":      users,
		"task":       task,
	}, nil)
}

func (h *TaskViewHandler) Edit(c *gin.Context) {
	var users []model.User
	var task model.Task
	db.DB.Find(&users)
	db.DB.Preload("User").First(&task, c.Param("id"))

	flash, _ := c.Get("flash")

	view.HTML(c, http.StatusOK, "tasks/edit.jet", gin.H{
		"title":      task.Title,
		"task":       task,
		"taskStatus": []string{"todo", "inprogress", "done"},
		"users":      users,
		"flash":      flash,
	}, nil)
}
