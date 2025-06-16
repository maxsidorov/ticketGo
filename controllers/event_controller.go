package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"log"
	"net/http"
	"sort"
	"time"
)

func MainPage(c *gin.Context) {
	butSearch := c.Query("but_search")
	search := c.Query("search")
	butSort := c.Query("but_sort")
	typeSort := c.Query("sort")
	dateFirst := c.Query("date_first")
	dateLast := c.Query("date_second")

	var events []models.Event

	query := DB
	if butSearch != "" || search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if dateFirst != "" {
		startDate, err := time.Parse("2006-01-02", dateFirst)
		handleError(err)
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(),
			0, 0, 0, 0,
			time.Local)
		query = query.Where("date_time >= ?", startDate)
	}
	if dateLast != "" {
		endDate, err := time.Parse("2006-01-02", dateLast)
		handleError(err)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day()+1,
			0, 0, 0, 0,
			time.Local)
		query = query.Where("date_time < ?", endDate)
	}
	query.Find(&events)
	if butSort != "" {
		switch typeSort {
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
		"date_first": dateFirst,
		"date_last":  dateLast,
		"username":   c.GetString("username"),
	})
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
