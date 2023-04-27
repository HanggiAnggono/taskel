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
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var errorMessage string
	var user model.User
	reqBody := LoginRequest{}
	c.ShouldBind(&reqBody)

	handleError := func() {
		switch c.Request.Header.Get("Content-Type") {
		case "application/json":
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": errorMessage,
			})
		default:
			c.HTML(http.StatusBadRequest, "auth/login", gin.H{
				"title": "Login",
				"error": errorMessage,
			})
		}
	}

	db.DB.Where("username = ?", reqBody.Username).First(&user)

	if user.ID == 0 {
		errorMessage = "user not found"
		handleError()
		return
	}

	if !model.UserComparePassword(reqBody.Password, user.Password) {
		errorMessage = "invalid password"
		handleError()
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Config.JWTSecret))

	if err != nil {
		errorMessage = fmt.Sprintf("failed to sign token: %v", err)
		handleError()
		return
	}

	switch c.Request.Header.Get("Content-Type") {
	case "application/json":
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "login successful",
			"user":    user,
			"token":   tokenString,
		})
	default:
		c.SetCookie("token", tokenString, 3600, "/", c.Request.Host, false, true)
		c.Redirect(http.StatusFound, "/")
	}
}

func (h *AuthHandler) Logout(c *gin.Context) {
	accept := c.Request.Header.Get("Accept")
	c.SetCookie("token", "", -1, "/", c.Request.Host, false, true)

	switch accept {
	case "application/json":
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "logout successful",
		})
	case "text/html":
	default:
		c.Redirect(http.StatusFound, "/login")
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
