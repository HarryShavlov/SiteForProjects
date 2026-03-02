package handlers

import (
	"net/http"

	"site/internal/models"
	"site/internal/services"
	"site/internal/store"

	"github.com/gin-gonic/gin"
)

// Login обрабатывает вход пользователя
func Login(c *gin.Context, ds *store.Store) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// Поиск пользователя
	user, found := ds.GetUserByLogin(req.Login)
	if !found || !services.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	// Генерация токена
	token, err := services.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	// Возвращаем пользователя без пароля
	userSafe := *user
	userSafe.Password = ""

	c.JSON(http.StatusOK, models.AuthResponse{Token: token, User: userSafe})
}

// GetMe возвращает данные текущего пользователя
func GetMe(c *gin.Context, ds *store.Store) {
	userID, _ := c.Get("userID")

	// Конвертируем interface{} в int
	var id int
	switch v := userID.(type) {
	case int:
		id = v
	case float64:
		id = int(v)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	user, found := ds.GetUserByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Не возвращаем хеш пароля
	userSafe := *user
	userSafe.Password = ""

	c.JSON(http.StatusOK, userSafe)
}
