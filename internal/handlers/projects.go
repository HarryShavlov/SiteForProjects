package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"site/internal/models"
	"site/internal/services"
	"site/internal/store"

	"github.com/gin-gonic/gin"
)

func GetProjects(c *gin.Context, ds *store.Store) {
	c.JSON(http.StatusOK, ds.GetProjects())
}

func CreateProject(c *gin.Context, ds *store.Store) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	// TODO: добавить в хранилище
	c.JSON(http.StatusCreated, project)
}

func ReloadProjects(c *gin.Context, ds *store.Store) {
	projects, err := services.LoadProjectsFromFile("data/projects.xlsx")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	users, err := services.LoadUsersFromFile("data/users.xlsx")
	if err != nil {
		users = getDemoUsers()
	}
	ds.SetProjects(projects)
	ds.SetUsers(users)

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"projects": len(projects),
		"users":    len(users),
	})
}

func UploadProjects(c *gin.Context, ds *store.Store) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не получен"})
		return
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Только файлы .xlsx"})
		return
	}

	tempPath := "data/temp_" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения"})
		return
	}

	newProjects, err := services.ParseProjectsFromTempFile(tempPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка чтения: " + err.Error()})
		return
	}

	ds.SetProjects(newProjects)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": fmt.Sprintf("Загружено %d проектов", len(newProjects)),
		"count":   len(newProjects),
	})
}

func getDemoUsers() []models.User {
	return []models.User{
		{ID: 1, Login: "admin", Password: services.MustHash("admin123"), Role: "admin", Name: "Администратор"},
		{ID: 2, Login: "student", Password: services.MustHash("student123"), Role: "student", Name: "Иван Иванов", Class: "8Б"},
	}
}
