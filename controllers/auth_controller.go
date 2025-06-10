package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/maxsidorov/ticketGo/models"
	"gorm.io/gorm"
)

var DB *gorm.DB // должен быть инициализирован в main.go

func ShowLoginPage(c *gin.Context) {
	session := sessions.Default(c)
	flash := session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "login.html", gin.H{"flash": flash})
}

func ShowRegisterPage(c *gin.Context) {
	session := sessions.Default(c)
	flash := session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "register.html", gin.H{"flash": flash})
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	var user models.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		session.AddFlash("Пользователь не найден")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func Register(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	if username == "" {
		session.AddFlash("Имя пользователя не может быть пустым")
		session.Save()
		c.Redirect(http.StatusFound, "/register")
		return
	}
	var count int64
	DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		session.AddFlash("Пользователь с таким именем уже существует")
		session.Save()
		c.Redirect(http.StatusFound, "/register")
		return
	}
	user := models.User{Username: username}
	if err := DB.Create(&user).Error; err != nil {
		session.AddFlash("Ошибка регистрации")
		session.Save()
		c.Redirect(http.StatusFound, "/register")
		return
	}
	session.AddFlash("Регистрация успешна! Войдите.")
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
} 