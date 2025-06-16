package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"net/http"
	"sort"
)

func MainPage(c *gin.Context) {
	but_search := c.Query("but_search")
	search := c.Query("search")
	but_sort := c.Query("but_sort")
	type_sort := c.Query("sort")

	var events []models.Event

	query := DB
	if but_search != "" || search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	query.Find(&events)
	if but_sort != "" {
		switch type_sort {
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
		"search":   search,
		"username": c.GetString("username"),
	})
}
