package service

import (
	"fmt"
	"taskel/config"
	"taskel/db"
	model "taskel/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthService struct{}

func (s *AuthService) GenerateJWTToken(userId uint, username string, roleId *uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"roleId":   roleId,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(config.Config.JWTSecret))
}

func (s *AuthService) ValidateJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWTSecret), nil
	})
}

func (s *AuthService) GetJWTClaims(tokenString string) (jwt.MapClaims, string, error) {
	token, err := s.ValidateJWTToken(tokenString)
	if err != nil {
		return nil, tokenString, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !token.Valid || !ok {
		return nil, tokenString, err
	}

	return claims, tokenString, nil
}

func (s *AuthService) GetJWTClaimsFromContext(c *gin.Context) (jwt.MapClaims, string, error) {
	token, _ := c.Cookie("token")
	return s.GetJWTClaims(token)
}

func (s *AuthService) IsAuthorized(c *gin.Context, permissionName string) bool {
	claims, _, _ := s.GetJWTClaimsFromContext(c)

	roleId := claims["roleId"]
	if roleId == nil {
		return false
	}

	fmt.Printf("roleId: %d\n", roleId)
	var permission model.Permission
	db.DB.Model(&model.Permission{}).Where("name = ?", permissionName).First(&permission)
	var count *int64

	db.DB.Table("role_permissions").
		Where("role_id = ?", roleId).
		Where("permission_id = ?", permission.ID).
		Count(count)

	return *count > 0
}
