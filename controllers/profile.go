package controllers

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/db"
)

func Profile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.Redirect(302, "/login")
		return
	}

	var userTickets []models.UserTicket
	if err := db.DB.Preload("Event").Where("user_id = ?", user.(models.User).ID).Find(&userTickets).Error; err != nil {
		c.HTML(500, "error.html", gin.H{
			"error": "Ошибка при загрузке билетов",
		})
		return
	}

	now := time.Now()
	var upcomingTickets, pastTickets []models.UserTicket

	for _, ticket := range userTickets {
		if ticket.Event.DateTime.After(now) {
			upcomingTickets = append(upcomingTickets, ticket)
		} else {
			pastTickets = append(pastTickets, ticket)
		}
	}

	c.HTML(200, "profile.html", gin.H{
		"User":           user,
		"upcomingTickets": upcomingTickets,
		"pastTickets":    pastTickets,
	})
} 