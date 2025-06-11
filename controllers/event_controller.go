package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"log"
	"net/http"
	"sort"
)

func ShowMainPage(c *gin.Context) {
	q := c.Query("q")
	var events []models.Event
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
	var events []models.Event
	query := DB
	query.Find(&events)
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
		case "По категории":
			log.Printf("Category!!!!!!!!!!!")
		}
	} else {
		log.Printf("but")
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"events":   events,
		"username": c.GetString("username"),
	})
}
