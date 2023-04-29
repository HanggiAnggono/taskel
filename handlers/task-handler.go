package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"taskel/db"
	model "taskel/models"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct{}

func (h *TaskHandler) List(c *gin.Context) {
	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, pageSizeErr := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageErr != nil || pageSizeErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "page or pageSize is invalid",
		})
	}

	var tasks []model.Task
	db.DB.Limit(pageSize).Offset((page - 1) * pageSize).Find(&tasks)

	c.JSON(http.StatusOK, gin.H{
		"data":    tasks,
		"status":  http.StatusOK,
		"message": "success",
	})
}

func (h *TaskHandler) Show(c *gin.Context) {
	id := c.Param("id")

	var task model.Task
	result := db.DB.First(&task, id)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "task not found",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"status":  http.StatusOK,
		"message": "success",
	})
}

type CreateRequest struct {
	Title       string  `json:"title" form:"title"`
	Status      string  `json:"status" form:"status"`
	UserID      *uint   `json:"userId" form:"userId"`
	Description *string `json:"description" form:"description"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateRequest
	c.ShouldBind(&req)

	var task model.Task
	task.Title = req.Title
	task.Status = req.Status
	task.Description = req.Description
	task.UserID = req.UserID

	db.DB.Save(&task)

	// handle Error
	if task.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "failed to create task",
		})
		return
	}

	content := c.Request.Header.Get("Content-Type")
	fmt.Printf("Content %s\n", content)
	switch content {
	case "application/json":
		c.JSON(http.StatusOK, gin.H{
			"data":    task,
			"status":  http.StatusOK,
			"message": "success",
		})
	case "application/x-www-form-urlencoded":
		c.Redirect(http.StatusMovedPermanently, "/")
	default:
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}

type EditRequest struct {
	Title       string  `json:"title" form:"title"`
	Status      string  `json:"status" form:"status"`
	Description *string `json:"description" form:"description"`
	UserID      *uint   `json:"userId" form:"userId"`
}

func (h *TaskHandler) Edit(c *gin.Context) {
	id := c.Param("id")
	req := EditRequest{}
	var task model.Task
	c.ShouldBind(&req)

	db.DB.First(&task, id)

	task.Title = req.Title
	task.Status = req.Status
	task.Description = req.Description
	task.UserID = req.UserID

	db.DB.Save(&task)

	switch c.Request.Header.Get("Content-Type") {
	case "application/json":
		c.JSON(http.StatusOK, gin.H{
			"data":    task,
			"status":  http.StatusOK,
			"message": "success",
		})
	case "application/x-www-form-urlencoded":
		taskViewHandler := TaskViewHandler{}
		c.Set("flash", "success")
		taskViewHandler.Edit(c)
	default:
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/task/%s/edit", id))
	}

}

type AssignRequest struct {
	UserID *uint `json:"userId"`
}

func (h *TaskHandler) AssignUserToTask(c *gin.Context) {
	id := c.Param("id")
	reqBody := AssignRequest{}
	c.BindJSON(&reqBody)
	userId := reqBody.UserID

	var task model.Task
	db.DB.First(&task, id)

	tx := db.DB.Begin()
	if userId != nil {
		task.UserID = userId
	}
	db.DB.Model(&task).UpdateColumn("user_id", userId)
	tx.Commit()

	if tx.Error != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "error",
			"message": tx.Error.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data":    task,
			"status":  http.StatusOK,
			"message": "success",
		})
	}
}

type TransitionRequest struct {
	Status string `json:"status"`
}

func (h *TaskHandler) TransitionTask(c *gin.Context) {
	id := c.Param("id")
	reqBody := TransitionRequest{}
	c.BindJSON(&reqBody)

	var task model.Task
	db.DB.First(&task, id)

	task.Status = reqBody.Status
	db.DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"status":  http.StatusOK,
		"message": "success",
	})
}