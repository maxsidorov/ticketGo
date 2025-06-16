package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"net/http"
	"sort"
	"time"
)

func MainPage(c *gin.Context) {
	but_search := c.Query("but_search")
	search := c.Query("search")
	but_sort := c.Query("but_sort")
	type_sort := c.Query("sort")
	date_first := c.Query("date_first")
	date_last := c.Query("date_second")

	var events []models.Event

	query := DB
	if but_search != "" || search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if date_first != "" {
		startDate, _ := time.Parse("2006-01-02", date_first)
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(),
			0, 0, 0, 0,
			time.Local)
		query = query.Where("date_time >= ?", startDate)
	}
	if date_last != "" {
		endDate, _ := time.Parse("2006-01-02", date_last)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day()+1,
			0, 0, 0, 0,
			time.Local)
		query = query.Where("date_time < ?", endDate)
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
		"events":     events,
		"search":     search,
		"date_first": date_first,
		"date_last":  date_last,
		"username":   c.GetString("username"),
	})
}
