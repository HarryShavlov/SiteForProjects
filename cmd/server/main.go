package main

import (
	"fmt"
	"log"
	"net/http"

	"site/internal/config"
	"site/internal/handlers"
	"site/internal/middleware"
	"site/internal/models"
	"site/internal/services"
	"site/internal/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация конфигурации
	config.Init()

	// Инициализация хранилища
	dataStore := store.New()

	// Загрузка данных при старте
	if err := loadData(dataStore); err != nil {
		log.Printf("⚠️ Ошибка загрузки данных: %v, используем демо-данные", err)
		loadDemoData(dataStore)
	}

	// Инициализация Gin
	r := setupRouter(dataStore)

	// Запуск сервера
	addr := fmt.Sprintf(":%s", config.App.ServerPort)
	fmt.Printf("🚀 Сервер запущен: http://localhost%s\n", config.App.ServerPort)

	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Ошибка запуска сервера: %v", err)
	}
}

func setupRouter(ds *store.Store) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	// Публичные маршруты
	r.POST("/api/login", func(c *gin.Context) {
		handlers.Login(c, ds)
	})
	r.GET("/api/projects", func(c *gin.Context) {
		handlers.GetProjects(c, ds)
	})

	// Защищённые маршруты
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", func(c *gin.Context) { handlers.GetMe(c, ds) })
		auth.POST("/projects", middleware.AdminOnly(), func(c *gin.Context) { handlers.CreateProject(c, ds) })
		auth.POST("/reload", middleware.AdminOnly(), func(c *gin.Context) { handlers.ReloadProjects(c, ds) })
		auth.POST("/upload", middleware.AdminOnly(), func(c *gin.Context) { handlers.UploadProjects(c, ds) })
	}

	// Статика
	r.Static("/static", config.App.StaticDir)
	r.StaticFile("/", config.App.StaticDir+"/index.html")

	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	return r
}

func loadData(ds *store.Store) error {
	projects, err := services.LoadProjectsFromFile("data/projects.xlsx")
	if err != nil {
		return err
	}
	users, err := services.LoadUsersFromFile("data/users.xlsx")
	if err != nil {
		return err
	}
	ds.SetProjects(projects)
	ds.SetUsers(users)
	return nil
}

func loadDemoData(ds *store.Store) {
	ds.SetUsers([]models.User{
		{ID: 1, Login: "admin", Password: services.MustHash("admin123"), Role: "admin", Name: "Администратор"},
		{ID: 2, Login: "student", Password: services.MustHash("student123"), Role: "student", Name: "Иван Иванов", Class: "8Б"},
	})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
