package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/config"
	"github.com/maxsidorov/ticketGo/controllers"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/routes"
	"github.com/maxsidorov/ticketGo/storage"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func main() {
	cfg := config.Load()
	var err error
	db, err = storage.InitPostgresGorm(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	// Миграция моделей
	db.AutoMigrate(&models.Event{}, &models.User{}, &models.Ticket{}, &models.UserTicket{})

	controllers.DB = db

	r := gin.Default()

	sessionStore := cookie.NewStore([]byte("secret-key-123"))
	r.Use(sessions.Sessions("session", sessionStore))

	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		username, _ := session.Get("username").(string)
		if username != "" {
			c.Set("username", username)
		}
		c.Next()
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")

	routes.RegisterRoutes(r)

	r.Run(":" + cfg.Port)
}
