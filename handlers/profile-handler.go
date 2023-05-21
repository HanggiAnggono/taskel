package handler

import (
	"fmt"
	"net/http"
	service "taskel/service"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct{}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	token, _ := c.Cookie("token")
	fmt.Printf("Token %s\n", token)
	authService := service.AuthService{}
	claims, _, err := authService.GetJWTClaims(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": claims,
	})
}
