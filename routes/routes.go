package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/controllers"
	"github.com/maxsidorov/ticketGo/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Статические файлы
	r.Static("/static", "./static")

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
	r.GET("/events/:id", controllers.ShowEventDetails)
	r.POST("/events/:id/buy", middleware.AuthRequired(), controllers.BuyTicket)

	// Обработка 404
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{
			"error": "Страница не найдена",
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
	})
}
