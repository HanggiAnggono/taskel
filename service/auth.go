package service

import (
	"taskel/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct{}

func (s *AuthService) GenerateJWTToken(userId uint, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": username,
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
