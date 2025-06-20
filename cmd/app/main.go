package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/config"
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/middleware"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/routes"
	"github.com/maxsidorov/ticketGo/storage"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"runtime/debug"
	"text/template"
	"time"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func main() {
	debug.SetGCPercent(200)
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

	r := gin.Default()

	sessionStore := cookie.NewStore([]byte("Mhvr56GCf65ucfjHD65fjwedo"))
	sessionStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 дней
		HttpOnly: true,
		Secure:   false, // true только если у тебя HTTPS
	})
	r.Use(sessions.Sessions("session", sessionStore))

	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		userID := session.Get("user_id")

		if username != nil {
			c.Set("username", username)
		}

		if userID != nil {
			c.Set("user_id", userID)
		}

		c.Next()
	})

	// Добавляем функции для шаблонов
	r.SetFuncMap(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("02.01.2006 15:04")
		},
		"sequence": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")

	// Установка middleware для всех запросов
	r.Use(middleware.SetAuthStatus())

	// Регистрируем все маршруты
	routes.RegisterRoutes(r, db.DB)
	// Создание пользователя admin, если его нет
	var adminUser models.User
	if err := db.DB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			password, err := hashPassword("admin")
			if err != nil {
				log.Printf("Не удалось захешировать пароль для admin: %v", err)
			} else {
				admin := models.User{
					Username:   "admin",
					Password:   password,
					Email:      "admin@admin.com",
					AdminLevel: 2,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				if err := db.DB.Create(&admin).Error; err != nil {
					log.Printf("Ошибка при создании пользователя admin: %v", err)
				} else {
					log.Printf("Пользователь admin успешно создан")
				}
			}
		}
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		panic(fmt.Sprintf("failed to start server: %v", err))
	}
}
