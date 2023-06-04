package handler

import (
	"fmt"
	"log"
	"net/http"
	"taskel/constants"
	"taskel/db"
	"taskel/mail_service"
	model "taskel/models"
	"taskel/repository"
	"taskel/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct{}

func (h *TaskHandler) List(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "error",
				"message": err,
			})
		}
	}()

	page, pageSize, _ := GetPaginationParams(c)

	var tasks []model.Task
	db.DB.Limit(pageSize).Offset((page - 1) * pageSize).Preload("User").Find(&tasks)

	c.JSON(http.StatusOK, gin.H{
		"data":    tasks,
		"status":  http.StatusOK,
		"message": "success",
	})
}

func (h *TaskHandler) Show(c *gin.Context) {
	key := c.Param("key")
	task, err := repository.GetTaskByIdOrKey(key)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
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
	Title       *string `json:"title" form:"title"`
	Status      *string `json:"status" form:"status"`
	Description *string `json:"description" form:"description"`
	UserID      *uint   `json:"userId" form:"userId"`
}

func (h *TaskHandler) Edit(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "error",
				"message": r,
			})
		}
	}()

	key := c.Param("key")
	req := EditRequest{}
	var task model.Task
	c.ShouldBind(&req)

	authService := service.AuthService{}
	isAuthorized := authService.IsAuthorized(c, constants.RBAC_Task_Write)
	if !isAuthorized {
		panic("role is unauthorized")
	}

	db.DB.Preload("Watchers").Where("key = ?", key).First(&task)

	// oldTask := task
	if (req.Title != nil) && (req.Title != &task.Title) {
		task.Title = *req.Title
	}
	if (req.Status != nil) && (req.Status != &task.Status) {
		task.Status = *req.Status
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if (req.UserID != nil) && (req.UserID != task.UserID) {
		task.UserID = req.UserID
	}

	db.DB.Save(&task)

	watcherEmails := []string{}
	if task.Watchers != nil {
		for _, watcher := range task.Watchers {
			watcherEmails = append(watcherEmails, *watcher.Email)
		}
	}

	if len(watcherEmails) > 0 {
		go mail_service.SendMail(
			fmt.Sprintf("There has been update on %v", task.Title),
			fmt.Sprintf("Title: %v\nStatus: %v\nDescription: %v", task.Title, task.Status, *task.Description),
			watcherEmails...,
		)
	}

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
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/task/%s/edit", key))
	}

}

type AssignRequest struct {
	UserID *uint `json:"userId"`
}

// @deprecated
func (h *TaskHandler) AssignUserToTask(c *gin.Context) {
	key := c.Param("key")
	reqBody := AssignRequest{}
	c.BindJSON(&reqBody)
	userId := reqBody.UserID

	var task model.Task
	db.DB.Where("key = ?", key).First(&task)

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

// Deprecated: use edit instead
func (h *TaskHandler) TransitionTask(c *gin.Context) {
	key := c.Param("key")
	reqBody := TransitionRequest{}
	c.BindJSON(&reqBody)

	var task model.Task
	db.DB.Where("key = ?", key).First(&task)

	task.Status = reqBody.Status
	db.DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"status":  http.StatusOK,
		"message": "success",
	})
}

type WatchTaskRequest struct {
	UserID uint `json:"userId" form:"userId"`
}

func (h *TaskHandler) WatchTask(c *gin.Context) {
	key := c.Param("key")
	reqBody := WatchTaskRequest{}
	c.ShouldBind(&reqBody)

	err := repository.TaskWatch(key, reqBody.UserID)
	taskPath := fmt.Sprintf("/task/%d", key)
	log.Printf("taskPath %s\n", taskPath)

	handleError := func() {
		switch c.Request.Header.Get("Content-Type") {
		case "application/json":
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		case "application/x-www-form-urlencoded":
			c.Redirect(http.StatusMovedPermanently, taskPath)
		}
	}

	if err != nil {
		handleError()
		return
	}

	switch c.Request.Header.Get("Content-Type") {
	case "application/json":
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success",
		})
	case "application/x-www-form-urlencoded":
		c.Redirect(http.StatusMovedPermanently, taskPath)
	}
}
