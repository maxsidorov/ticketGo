package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/service"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// checkPasswordHash проверяет соответствие хеша пароля
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// hashPassword создает хеш пароля
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func Login(c *gin.Context) {
	usernameOrEmail := c.PostForm("username")
	password := c.PostForm("password")

	// Ищем пользователя по имени пользователя или email
	var user models.User
	if err := db.DB.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error":         "Неверное имя пользователя/email или пароль",
			"IsAuthenticated": false,
		})
		return
	}

	// Проверяем пароль
	if !checkPasswordHash(password, user.Password) {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error":         "Неверное имя пользователя/email или пароль",
			"IsAuthenticated": false,
		})
		return
	}

	// Сохраняем информацию в сессии
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error":         "Ошибка при сохранении сессии",
			"IsAuthenticated": false,
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Валидация имени пользователя
	if err, _ := service.ValidateName(username); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Валидация пароля
	if err, _ := service.ValidatePassword(password); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Проверяем совпадение паролей
	if password != confirmPassword {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Пароли не совпадают",
		})
		return
	}

	// Проверяем, существует ли пользователь с таким именем
	var existingUser models.User
	if err := db.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Пользователь с таким именем уже существует",
		})
		return
	}

	// Проверяем, существует ли пользователь с таким email
	if err := db.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Пользователь с таким email уже существует",
		})
		return
	}

	// Создаем хеш пароля
	hashedPassword, err := hashPassword(password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Ошибка при создании пользователя",
		})
		return
	}

	// Создаем нового пользователя
	user := models.User{
		Username:   username,
		Email:      email,
		Password:   hashedPassword,
		AdminLevel: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Ошибка при создании пользователя",
		})
		return
	}

	// Автоматически входим в систему после регистрации
	session := sessions.Default(c)
	session.Set("username", username)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Ошибка при сохранении сессии",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func CheckAuth(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	c.JSON(http.StatusOK, gin.H{
		"authenticated": userID != nil,
	})
}
