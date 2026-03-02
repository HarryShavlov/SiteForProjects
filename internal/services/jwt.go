package services

import (
	"errors"
	"time"

	"site/internal/config"
	"site/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken создаёт JWT токен для пользователя
func GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"id":    user.ID,
		"login": user.Login,
		"role":  user.Role,
		"exp":   expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.App.JWTSecret)
}

// ParseToken проверяет и парсит JWT токен
func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.App.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("неверный токен")
	}
	return claims, nil
}
