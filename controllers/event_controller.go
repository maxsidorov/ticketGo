package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"net/http"
	"sort"
)

var events []models.Event

func ShowMainPage(c *gin.Context) {
	q := c.Query("q")
	query := DB
	if q != "" {
		query = query.Where("title ILIKE ?", "%"+q+"%")
	}
	query.Find(&events)
	// TODO: username из сессии
	c.HTML(http.StatusOK, "index.html", gin.H{
		"events":   events,
		"username": c.GetString("username"),
	})
}

func MainPage(c *gin.Context) {
	p := c.PostForm("but")
	r := c.PostForm("phone")
	if p != "" {
		switch r {
		case "SortName":
			sort.Slice(events, func(i, j int) bool { return events[i].Title < events[j].Title })
		case "SortPrice":
			sort.Slice(events, func(i, j int) bool { return events[i].Price < events[j].Price })
		case "SortDate":
			sort.Slice(events, func(i, j int) bool { return events[i].DateTime.Before(events[j].DateTime) })
		case "SortCategory":
			events = events[1:]
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"events":   events,
		"username": c.GetString("username"),
	})
}
