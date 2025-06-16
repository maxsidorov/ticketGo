package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var DB *gorm.DB // должен быть инициализирован в main.go

func ShowLoginPage(c *gin.Context) {
	session := sessions.Default(c)
	flash := session.Flashes()
	sessionSave(session)
	c.HTML(http.StatusOK, "login.html", gin.H{"flash": flash})
}

func ShowRegisterPage(c *gin.Context) {
	session := sessions.Default(c)
	flash := session.Flashes()
	sessionSave(session)
	c.HTML(http.StatusOK, "register.html", gin.H{"flash": flash})
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	userpass := c.PostForm("userpass")
	var user models.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		session.AddFlash("Пользователь не найден")
		sessionSave(session)
		c.Redirect(http.StatusFound, "/login")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userpass)); err != nil {
		session.AddFlash("Неверный пароль")
		sessionSave(session)
		c.Redirect(http.StatusFound, "/login")
		return
	}
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	sessionSave(session)
	c.Redirect(http.StatusFound, "/")
}

func Register(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	userpass := c.PostForm("userpass")

	errName, username := service.ValidateName(username)
	if errName != nil {
		session.AddFlash(errName.Error())
		sessionSave(session)
		c.Redirect(http.StatusFound, "/register")
		return
	}
	var count int64
	DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		session.AddFlash("Пользователь с таким именем уже существует")
		sessionSave(session)
		c.Redirect(http.StatusFound, "/register")
		return
	}

	errPass, userpass := service.ValidatePassword(userpass)
	if errPass != nil {
		session.AddFlash(errPass.Error())
		sessionSave(session)
		c.Redirect(http.StatusFound, "/register")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(userpass), bcrypt.DefaultCost)
	if err != nil {
		session.AddFlash("Ошибка при обработке пароля")
		sessionSave(session)
		c.Redirect(http.StatusFound, "/register")
		return
	}
	user := models.User{Username: username, Password: string(hash), IsAdmin: false}
	if err := DB.Create(&user).Error; err != nil {
		session.AddFlash("Ошибка регистрации")
		sessionSave(session)
		c.Redirect(http.StatusFound, "/register")
		return
	}
	session.AddFlash("Регистрация успешна! Войдите.")
	sessionSave(session)
	c.Redirect(http.StatusFound, "/login")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	sessionSave(session)
	c.Redirect(http.StatusFound, "/")
}

func sessionSave(session sessions.Session) {
	err := session.Save()
	if err != nil {
		log.Print("Session error: ", err)
	}

}
