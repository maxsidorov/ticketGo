package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/validator"
)

func ShowProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := db.DB.Preload("Tickets.Event").First(&user, userID).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"error":         "Пользователь не найден",
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
		return
	}

	// Разделяем билеты на предстоящие и прошедшие
	var upcomingTickets, pastTickets []models.UserTicket
	now := time.Now()

	for _, ticket := range user.Tickets {
		if ticket.Event.DateTime.After(now) {
			upcomingTickets = append(upcomingTickets, ticket)
		} else {
			pastTickets = append(pastTickets, ticket)
		}
	}

	// Сортируем билеты по дате
	sortTicketsByDate := func(tickets []models.UserTicket) {
		for i := 0; i < len(tickets)-1; i++ {
			for j := i + 1; j < len(tickets); j++ {
				if tickets[i].Event.DateTime.After(tickets[j].Event.DateTime) {
					tickets[i], tickets[j] = tickets[j], tickets[i]
				}
			}
		}
	}

	sortTicketsByDate(upcomingTickets)
	sortTicketsByDate(pastTickets)

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"User":            user,
		"UpcomingTickets": upcomingTickets,
		"PastTickets":     pastTickets,
		"IsAuthenticated": true,
		"Username":        user.Username,
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
		Username        string `json:"username"`
		Email          string `json:"email"`
		NewPassword    string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Неверный формат данных",
		})
		return
	}

	// Проверяем, не занято ли имя пользователя
	if formData.Username != user.Username {
		var existingUser models.User
		if err := db.DB.Where("username = ?", formData.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "username already exists",
			})
			return
		}
		user.Username = formData.Username
	}

	// Проверяем, не занят ли email
	if formData.Email != user.Email {
		var existingUser models.User
		if err := db.DB.Where("email = ?", formData.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "email already exists",
			})
			return
		}
		user.Email = formData.Email
	}

	// Если пользователь хочет изменить пароль
	if formData.NewPassword != "" {
		if err := validator.ValidatePassword(formData.NewPassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid password format",
			})
			return
		}

		// Хешируем новый пароль
		hashedPassword, err := validator.HashPassword(formData.NewPassword)
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Настройки успешно сохранены",
	})
} 