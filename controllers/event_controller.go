package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/db"
	"net/http"
	"strconv"
	"gorm.io/gorm"
	"github.com/gin-contrib/sessions"
	"log"
	"time"
	"math"
	"fmt"
)

type EventController struct {
	db *gorm.DB
}

func NewEventController(db *gorm.DB) *EventController {
	return &EventController{
		db: db,
	}
}

func ShowMainPage(c *gin.Context) {
	var events []models.Event
	searchQuery := c.Query("search")
	sortType := c.Query("sort")
	category := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "12")
	
	log.Printf("Starting ShowMainPage with search: %s, sort: %s, category: %s, minPrice: %s, maxPrice: %s, startDate: %s, endDate: %s, page: %s", 
		searchQuery, sortType, category, minPrice, maxPrice, startDate, endDate, page)
	
	// Преобразуем параметры пагинации
	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSizeNum < 1 || pageSizeNum > 100 {
		pageSizeNum = 12
	}
	offset := (pageNum - 1) * pageSizeNum
	
	// Создаем базовый запрос
	query := db.DB.Model(&models.Event{})
	
	// Применяем поиск с использованием ILIKE для регистронезависимого поиска
	if searchQuery != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", 
			"%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Применяем фильтр по категории
	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	// Применяем фильтр по цене
	if minPrice != "" {
		if minPriceFloat, err := strconv.ParseFloat(minPrice, 64); err == nil {
			log.Printf("Applying min price filter: %f", minPriceFloat)
			query = query.Where("price >= ?", minPriceFloat)
		} else {
			log.Printf("Error parsing min price: %v", err)
		}
	}
	if maxPrice != "" {
		if maxPriceFloat, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			log.Printf("Applying max price filter: %f", maxPriceFloat)
			query = query.Where("price <= ?", maxPriceFloat)
		} else {
			log.Printf("Error parsing max price: %v", err)
		}
	}

	// Применяем фильтр по дате
	if startDate != "" {
		if startDateTime, err := time.Parse("2006-01-02", startDate); err == nil {
			log.Printf("Applying start date filter: %v", startDateTime)
			query = query.Where("date_time >= ?", startDateTime)
		} else {
			log.Printf("Error parsing start date: %v", err)
		}
	}
	if endDate != "" {
		if endDateTime, err := time.Parse("2006-01-02", endDate); err == nil {
			// Добавляем 23:59:59 к дате окончания для включения всего дня
			endDateTime = endDateTime.Add(24*time.Hour - time.Second)
			log.Printf("Applying end date filter: %v", endDateTime)
			query = query.Where("date_time <= ?", endDateTime)
		} else {
			log.Printf("Error parsing end date: %v", err)
		}
	}
	
	// Получаем текущее время
	now := time.Now()
	
	// Разделяем запрос на предстоящие и прошедшие события
	var upcomingEvents, pastEvents []models.Event
	
	// Получаем предстоящие события
	upcomingQuery := query.Where("date_time >= ?", now)
	switch sortType {
	case "date":
		upcomingQuery = upcomingQuery.Order("date_time ASC")
	case "date-desc":
		upcomingQuery = upcomingQuery.Order("date_time DESC")
	case "price-asc":
		upcomingQuery = upcomingQuery.Order("price ASC")
	case "price-desc":
		upcomingQuery = upcomingQuery.Order("price DESC")
	case "title":
		upcomingQuery = upcomingQuery.Order("title ASC")
	case "popular":
		upcomingQuery = upcomingQuery.Order("sold_tickets DESC")
	default:
		upcomingQuery = upcomingQuery.Order("date_time ASC")
	}
	
	if err := upcomingQuery.Find(&upcomingEvents).Error; err != nil {
		log.Printf("Error fetching upcoming events: %v", err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "Ошибка при получении мероприятий",
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
		return
	}
	
	// Получаем прошедшие события
	pastQuery := query.Where("date_time < ?", now).Order("date_time DESC")
	if err := pastQuery.Find(&pastEvents).Error; err != nil {
		log.Printf("Error fetching past events: %v", err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "Ошибка при получении мероприятий",
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
		return
	}
	
	// Объединяем события
	events = append(upcomingEvents, pastEvents...)
	
	// Применяем пагинацию
	total := len(events)
	start := offset
	end := offset + pageSizeNum
	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}
	
	// Получаем страницу событий
	pageEvents := events[start:end]
	
	// Получаем информацию о пользователе из сессии
	session := sessions.Default(c)
	username := session.Get("username")
	userID := session.Get("user_id")

	// Рассчитываем информацию о пагинации
	totalPages := int(math.Ceil(float64(total) / float64(pageSizeNum)))
	
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Events": pageEvents,
		"SearchQuery": searchQuery,
		"SortType": sortType,
		"Category": category,
		"MinPrice": minPrice,
		"MaxPrice": maxPrice,
		"StartDate": startDate,
		"EndDate": endDate,
		"IsAuthenticated": userID != nil,
		"Username": username,
		"Pagination": gin.H{
			"CurrentPage": pageNum,
			"TotalPages": totalPages,
			"TotalItems": total,
			"PageSize": pageSizeNum,
		},
	})
}

func (ec *EventController) ShowEvent(c *gin.Context) {
	eventID := c.Param("id")
	log.Printf("Loading event with ID: %s", eventID)

	var event models.Event
	if err := ec.db.First(&event, eventID).Error; err != nil {
		log.Printf("Error loading event: %v", err)
		panic("Событие не найдено")
	}

	log.Printf("Successfully loaded event: %+v", event)

	// Get user ID from session
	session := sessions.Default(c)
	userID := session.Get("user_id")
	username := session.Get("username")
	isAuthenticated := userID != nil

	var userTickets int
	if isAuthenticated {
		// Get user's tickets for this event
		var ticket models.UserTicket
		if err := ec.db.Where("user_id = ? AND event_id = ?", userID, eventID).First(&ticket).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Printf("Error checking user tickets: %v", err)
			}
			userTickets = 0
		} else {
			userTickets = ticket.Quantity
			log.Printf("Found user tickets: %d", userTickets)
		}
	}

	// Calculate remaining tickets
	remainingTickets := event.TotalTickets - event.SoldTickets
	log.Printf("Remaining tickets: %d", remainingTickets)

	c.HTML(http.StatusOK, "event.html", gin.H{
		"Event":           event,
		"UserTickets":     userTickets,
		"RemainingTickets": remainingTickets,
		"IsAuthenticated": isAuthenticated,
		"Username":        username,
	})
}

func (ec *EventController) BuyTicket(c *gin.Context) {
	eventID := c.Param("id")
	userID := sessions.Default(c).Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходимо авторизоваться"})
		return
	}

	// Получаем количество билетов из формы
	quantityStr := c.PostForm("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное количество билетов"})
		return
	}

	// Начинаем транзакцию
	tx := ec.db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке запроса"})
		return
	}

	// Получаем информацию о событии
	var event models.Event
	if err := tx.First(&event, eventID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Событие не найдено"})
		return
	}

	// Проверяем доступность билетов
	if event.SoldTickets+quantity > event.TotalTickets {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недостаточно билетов"})
		return
	}

	// Проверяем существующие билеты пользователя
	var existingTicket models.UserTicket
	err = tx.Where("user_id = ? AND event_id = ?", userID, eventID).First(&existingTicket).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Создаем новую запись о билетах
			newTicket := models.UserTicket{
				UserID:  userID.(uint),
				EventID: event.ID,
				Quantity: quantity,
			}
			if err := tx.Create(&newTicket).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при покупке билета"})
				return
			}
		} else {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке билетов"})
			return
		}
	} else {
		// Обновляем существующую запись
		existingTicket.Quantity += quantity
		if err := tx.Save(&existingTicket).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении билетов"})
			return
		}
	}

	// Обновляем количество проданных билетов
	event.SoldTickets += quantity
	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении события"})
		return
	}

	// Завершаем транзакцию
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при завершении транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Успешно куплено %d билетов", quantity),
		"tickets": quantity,
	})
}

func ShowEvents(c *gin.Context) {
	// Получаем параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Получаем параметры фильтрации
	category := c.Query("category")
	search := c.Query("search")
	sort := c.DefaultQuery("sort", "date")

	// Создаем базовый запрос
	query := db.DB.Model(&models.Event{})

	// Применяем фильтры
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Применяем сортировку
	switch sort {
	case "date":
		query = query.Order("date_time ASC")
	case "date-desc":
		query = query.Order("date_time DESC")
	case "price-asc":
		query = query.Order("price ASC")
	case "price-desc":
		query = query.Order("price DESC")
	case "title":
		query = query.Order("title ASC")
	case "popular":
		query = query.Order("sold_tickets DESC")
	default:
		query = query.Order("date_time ASC")
	}

	// Получаем общее количество событий
	var total int64
	query.Count(&total)

	// Применяем пагинацию
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Получаем события
	var events []models.Event
	if err := query.Find(&events).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при получении списка мероприятий",
		})
		return
	}

	// Рассчитываем информацию о пагинации
	totalPages := (int(total) + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "events.html", gin.H{
		"events":      events,
		"page":        page,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
		"total":       total,
		"category":    category,
		"search":      search,
		"sort":        sort,
		"categories":  []string{"concert", "theater", "exhibition", "sport", "other"},
	})
}

func (ec *EventController) ReturnTicket(c *gin.Context) {
	eventID := c.Param("id")
	userID := sessions.Default(c).Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходимо авторизоваться"})
		return
	}

	// Получаем количество билетов для возврата из формы
	quantityStr := c.PostForm("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное количество билетов"})
		return
	}

	// Начинаем транзакцию
	tx := ec.db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке запроса"})
		return
	}

	// Получаем информацию о событии
	var event models.Event
	if err := tx.First(&event, eventID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Событие не найдено"})
		return
	}

	// Проверяем существующие билеты пользователя
	var userTicket models.UserTicket
	err = tx.Where("user_id = ? AND event_id = ?", userID, eventID).First(&userTicket).Error

	if err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "У вас нет билетов на это мероприятие"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке билетов"})
		}
		return
	}

	// Проверяем, что пользователь не пытается вернуть больше билетов, чем у него есть
	if quantity > userTicket.Quantity {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя вернуть больше билетов, чем у вас есть"})
		return
	}

	// Проверяем, что мероприятие еще не прошло
	if event.DateTime.Before(time.Now()) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя вернуть билеты на прошедшее мероприятие"})
		return
	}

	// Обновляем количество билетов пользователя
	if quantity == userTicket.Quantity {
		// Если возвращаем все билеты, удаляем запись
		if err := tx.Delete(&userTicket).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении билетов"})
			return
		}
	} else {
		// Иначе уменьшаем количество
		userTicket.Quantity -= quantity
		if err := tx.Save(&userTicket).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении билетов"})
			return
		}
	}

	// Обновляем количество проданных билетов
	event.SoldTickets -= quantity
	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении события"})
		return
	}

	// Завершаем транзакцию
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при завершении транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Успешно возвращено %d билетов", quantity),
		"returned": quantity,
	})
}
