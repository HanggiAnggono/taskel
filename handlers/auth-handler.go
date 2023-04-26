package handler

import (
	"net/http"
	"taskel/db"
	model "taskel/models"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	reqBody := LoginRequest{}
	c.Bind(&reqBody)

	var user model.User
	db.DB.Where("username = ?", reqBody.Username).First(&user)

	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "user not found",
		})
	} else {
		if model.UserComparePassword(reqBody.Password, user.Password) {
			c.JSON(http.StatusOK, gin.H{
				"data":    user,
				"status":  http.StatusOK,
				"message": "success",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid password",
			})
		}
	}
}

func (h *AuthHandler) LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login", gin.H{
		"title": "Login",
	})
}
