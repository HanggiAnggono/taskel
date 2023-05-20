package handler

import (
	"net/http"
	"strconv"
	"taskel/db"
	model "taskel/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {}

func (h *UserHandler) List(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err,
			})
		}
	}()

	users := []model.User{}

	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, pageSizeErr := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageErr != nil || pageSizeErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "page or pageSize is invalid",
		})
	}

	db.DB.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}