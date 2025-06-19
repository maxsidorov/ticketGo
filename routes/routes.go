package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/controllers"
	"github.com/maxsidorov/ticketGo/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// Инициализация контроллеров
	eventController := controllers.NewEventController(db)

	// Главная страница
	r.GET("/", controllers.ShowMainPage)

	// Аутентификация
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)
	r.GET("/register", controllers.ShowRegisterPage)
	r.POST("/register", controllers.Register)
	r.GET("/logout", controllers.Logout)

	// Маршруты для профиля
	r.GET("/profile", middleware.AuthRequired(), controllers.ShowProfile)
	r.POST("/profile/update", middleware.AuthRequired(), controllers.UpdateProfile)

	// Маршруты для мероприятий
	r.GET("/events", controllers.ShowEvents)
	r.GET("/events/:id", eventController.ShowEvent)
	r.POST("/events/:id/buy", middleware.AuthRequired(), eventController.BuyTicket)
	r.GET("/admin", controllers.AdminPage)
	r.GET("/admin/events/new", controllers.AddEventPage)
	r.POST("/admin/events/new", controllers.AddEventPage)
	// Обработка 404
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{
			"title": "404 - Страница не найдена",
		})
	})
}
