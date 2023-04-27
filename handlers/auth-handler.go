package handler

import (
	"fmt"
	"net/http"
	"taskel/config"
	"taskel/db"
	model "taskel/models"
	"time"

	"github.com/dgrijalva/jwt-go"
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
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id":  user.ID,
				"username": user.Username,
				"exp":      time.Now().Add(time.Hour * 72).Unix(),
			})

			tokenString, err := token.SignedString([]byte(config.Config.JWTSecret))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("failed to sign token: %v", err),
				})
			} else {
				// success
				c.JSON(http.StatusOK, gin.H{
					"data": gin.H{
						"user":  user,
						"token": tokenString,
					},
					"status":  http.StatusOK,
					"message": "success",
				})
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid password",
			})
		}
	}
}

func (h *AuthHandler) LoginView(c *gin.Context) {
	jwtToken, _ := c.Cookie("token")

	if jwtToken != "" {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	c.HTML(http.StatusOK, "auth/login", gin.H{
		"title": "Login",
	})
}
