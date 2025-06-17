package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/db"
)

func GetEvent(c *gin.Context) {
	var event models.Event
	if err := db.DB.First(&event, c.Param("id")).Error; err != nil {
		c.HTML(404, "error.html", gin.H{
			"error": "Мероприятие не найдено",
		})
		return
	}

	user, exists := c.Get("user")
	var userTickets int
	if exists {
		var userTicket models.UserTicket
		if err := db.DB.Where("user_id = ? AND event_id = ?", user.(models.User).ID, event.ID).First(&userTicket).Error; err == nil {
			userTickets = userTicket.TicketsCount
		}
	}

	c.HTML(200, "event.html", gin.H{
		"Event":      event,
		"UserTickets": userTickets,
		"IsAuthenticated": exists,
	})
} 