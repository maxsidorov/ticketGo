package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/controllers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", controllers.ShowMainPage)
	r.POST("/", controllers.MainPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)
	r.GET("/register", controllers.ShowRegisterPage)
	r.POST("/register", controllers.Register)
	r.GET("/logout", controllers.Logout)
}
