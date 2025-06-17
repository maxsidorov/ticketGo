package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/service"
	"github.com/gin-contrib/sessions"
)

func ShowProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при получении данных пользователя",
		})
		return
	}

	var userTickets []models.UserTicket
	if err := db.DB.Preload("Event").Where("user_id = ?", userID).Find(&userTickets).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при получении билетов",
		})
		return
	}

	var upcomingTickets, pastTickets []models.UserTicket
	now := time.Now()
	for _, ut := range userTickets {
		if ut.Event.DateTime.After(now) {
			upcomingTickets = append(upcomingTickets, ut)
		} else {
			pastTickets = append(pastTickets, ut)
		}
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"Username":        username,
		"Email":           user.Email,
		"UpcomingTickets": upcomingTickets,
		"PastTickets":     pastTickets,
		"IsAuthenticated": true,
	})
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Пользователь не найден",
		})
		return
	}

	var formData struct {
		Username        string `json:"username" form:"username"`
		Email           string `json:"email" form:"email"`
		CurrentPassword string `json:"current_password" form:"current_password"`
		NewPassword     string `json:"new_password" form:"new_password"`
	}

	// Пробуем сначала получить данные из JSON
	if err := c.ShouldBindJSON(&formData); err != nil {
		// Если не получилось, пробуем получить из формы
		if err := c.ShouldBind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Неверный формат данных",
		})
		return
		}
	}

	// Проверяем текущий пароль, если он указан
	if formData.CurrentPassword != "" {
		if !service.CheckPasswordHash(formData.CurrentPassword, user.Password) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Неверный текущий пароль",
			})
			return
		}
	}

	// Проверяем, не занято ли имя пользователя
	if formData.Username != "" && formData.Username != user.Username {
		var existingUser models.User
		if err := db.DB.Where("username = ?", formData.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Имя пользователя уже занято",
			})
			return
		}
		user.Username = formData.Username
	}

	// Проверяем, не занят ли email
	if formData.Email != "" && formData.Email != user.Email {
		var existingUser models.User
		result := db.DB.Where("email = ? AND id != ?", formData.Email, user.ID).First(&existingUser)
		if result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Email уже занят",
			})
			return
		}
		user.Email = formData.Email
	}

	// Если пользователь хочет изменить пароль
	if formData.NewPassword != "" {
		if err, _ := service.ValidatePassword(formData.NewPassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// Хешируем новый пароль
		hashedPassword, err := service.HashPassword(formData.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Ошибка при обновлении пароля",
			})
			return
		}
		user.Password = hashedPassword
	}

	// Сохраняем изменения
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Ошибка при сохранении изменений",
		})
		return
	}

	// Обновляем сессию с новым именем пользователя
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Настройки успешно сохранены",
	})
} 