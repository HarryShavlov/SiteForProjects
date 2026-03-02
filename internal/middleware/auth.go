package middleware

import (
	"net/http"
	"strings"

	"site/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := services.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		// Сохраняем данные пользователя в контекст запроса
		c.Set("userID", (*claims)["id"])
		c.Set("userRole", (*claims)["role"])
		c.Next()
	}
}

// AdminOnly разрешает доступ только пользователям с ролью admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ только для администраторов"})
			c.Abort()
			return
		}
		c.Next()
	}
}
