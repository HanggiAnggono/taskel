package handler

import (
	"net/http"
	"taskel/db"
	model "taskel/models"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct{}

type CommentListParams struct {
	CommentableID   uint   `form:"commentable_id"`
	CommentableType string `form:"commentable_type"`
}

func (h *CommentHandler) List(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "error",
				"message": err,
			})
		}
	}()

	page, pageSize, _ := GetPaginationParams(c)
	params := CommentListParams{}
	c.ShouldBindQuery(&params)

	var comments []model.Comment

	db.DB.Limit(pageSize).Offset((page-1)*pageSize).Preload("Author").Where("commentable_id = ? AND commentable_type = ?", params.CommentableID, params.CommentableType).Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"data":    comments,
		"status":  http.StatusOK,
		"message": "success",
	})
}
