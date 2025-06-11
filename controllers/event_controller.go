package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"net/http"
	"sort"
)

// var DB *gorm.DB // должен быть инициализирован в main.go

func ShowMainPage(c *gin.Context) {
	q := c.Query("q")
	var events []models.Event
	query := DB

	if q != "" {
		query = query.Where("title ILIKE ?", "%"+q+"%")
	}
	query.Find(&events)
	p := c.Query("but")
	if p != "" {
		sort.Slice(events, func(i, j int) bool { return events[i].Title < events[j].Title })
	}
	// TODO: username из сессии
	c.HTML(http.StatusOK, "index.html", gin.H{
		"events":   events,
		"username": c.GetString("username"),
	})
}
