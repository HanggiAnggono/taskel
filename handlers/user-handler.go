package handler

import (
	"net/http"
	"taskel/db"
	model "taskel/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

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
	page, pageSize, _ := GetPaginationParams(c)

	db.DB.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}
