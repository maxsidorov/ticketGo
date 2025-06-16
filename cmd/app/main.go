package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/config"
	"github.com/maxsidorov/ticketGo/controllers"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/storage"
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"log"
	"text/template"
	"github.com/maxsidorov/ticketGo/middleware"
)

func main() {
	cfg := config.Load()
	var err error
	db.DB, err = storage.InitPostgresGorm(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %v", err))
	}

	log.Printf("AutoMigrate started")
	// Миграция моделей
	if err := db.DB.AutoMigrate(&models.Event{}, &models.User{}, &models.Ticket{}, &models.UserTicket{}); err != nil {
		panic(fmt.Sprintf("migration failed: %v", err))
	} else {
		log.Printf("AutoMigrate completed successfully")
	}

	controllers.DB = db.DB

	r := gin.Default()

	sessionStore := cookie.NewStore([]byte("secret-key-123"))
	r.Use(sessions.Sessions("session", sessionStore))

	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		userID := session.Get("user_id")
		
		if username != nil {
			c.Set("username", username)
		}
		
		if userID != nil {
			// Преобразуем ID в uint
			if userIDInt, ok := userID.(int); ok {
				c.Set("user_id", uint(userIDInt))
			}
		}
		
		c.Next()
	})

	// Добавляем функцию sub для шаблонов
	r.SetFuncMap(template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")

	// Установка middleware для всех запросов
	r.Use(middleware.SetAuthStatus())

	// Регистрируем все маршруты
	routes.RegisterRoutes(r)

	if err := r.Run(":" + cfg.Port); err != nil {
		panic(fmt.Sprintf("failed to start server: %v", err))
	}
}
