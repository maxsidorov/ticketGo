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

var DB *gorm.DB // должен быть инициализирован в main.go

var events []models.Event

func ShowMainPage(c *gin.Context) {
	var events []models.Event
	searchQuery := c.Query("search")
	sortType := c.Query("sort")
	category := c.Query("category")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "12")
	
	log.Printf("Starting ShowMainPage with search: %s, sort: %s, category: %s, page: %s", 
		searchQuery, sortType, category, page)
	
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
	
	// Применяем сортировку с учетом индексов
	switch sortType {
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
	
	// Получаем общее количество записей для пагинации
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting events: %v", err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "Ошибка при получении количества мероприятий",
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
		return
	}
	
	// Применяем пагинацию
	query = query.Offset(offset).Limit(pageSizeNum)
	
	// Выполняем запрос
	if err := query.Find(&events).Error; err != nil {
		log.Printf("Error fetching events: %v", err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "Ошибка при получении мероприятий",
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
		return
	}

	log.Printf("Found %d events", len(events))

	// Получаем информацию о пользователе из сессии
	session := sessions.Default(c)
	username := session.Get("username")
	userID := session.Get("user_id")

	// Рассчитываем информацию о пагинации
	totalPages := int(math.Ceil(float64(total) / float64(pageSizeNum)))
	
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Events": events,
		"SearchQuery": searchQuery,
		"SortType": sortType,
		"Category": category,
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

func MainPage(c *gin.Context) {
	var events []models.Event
	DB.Find(&events)

	// Получаем параметр сортировки
	sortType := c.PostForm("phone")

	// Применяем сортировку
	switch sortType {
	case "SortName":
		// Сортировка по названию
		for i := 0; i < len(events)-1; i++ {
			for j := i + 1; j < len(events); j++ {
				if events[i].Title > events[j].Title {
					events[i], events[j] = events[j], events[i]
				}
			}
		}
	case "SortPrice":
		// Сортировка по цене
		for i := 0; i < len(events)-1; i++ {
			for j := i + 1; j < len(events); j++ {
				if events[i].Price > events[j].Price {
					events[i], events[j] = events[j], events[i]
				}
			}
		}
	case "SortDate":
		// Сортировка по дате
		for i := 0; i < len(events)-1; i++ {
			for j := i + 1; j < len(events); j++ {
				if events[i].DateTime.After(events[j].DateTime) {
					events[i], events[j] = events[j], events[i]
				}
			}
		}
	case "SortCategory":
		// Заглушка для сортировки по категории - просто возвращаем список без изменений
		break
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"events":   events,
		"username": c.GetString("username"),
	})
}

// ShowEventDetails отображает страницу мероприятия
func ShowEventDetails(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid event ID: %v", err)
		c.HTML(http.StatusBadRequest, "404.html", gin.H{
			"error": "Неверный ID мероприятия",
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
		return
	}

	var event models.Event
	if err := db.DB.First(&event, eventID).Error; err != nil {
		log.Printf("Error fetching event: %v", err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"error": "Мероприятие не найдено",
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
		return
	}

	log.Printf("Found event: %+v", event)

	// Получаем количество проданных билетов
	var soldTickets int64
	db.DB.Model(&models.Ticket{}).Where("event_id = ?", event.ID).Count(&soldTickets)

	// Проверяем, авторизован ли пользователь
	session := sessions.Default(c)
	userID := session.Get("user_id")
	username := session.Get("username")
	isAuthenticated := userID != nil

	var userTicket models.UserTicket
	var hasTicket bool
	if isAuthenticated {
		if err := db.DB.Where("user_id = ? AND event_id = ?", userID, event.ID).First(&userTicket).Error; err == nil {
			hasTicket = true
		}
	}

	// Форматируем дату и цену для отображения
	formattedDate := event.DateTime.Format("02.01.2006 15:04")
	formattedPrice := fmt.Sprintf("%.2f", event.Price)

	data := gin.H{
		"Event": gin.H{
			"ID": event.ID,
			"Title": event.Title,
			"Image": event.Image,
			"DateTime": formattedDate,
			"Location": event.Location,
			"Price": formattedPrice,
			"Description": event.Description,
			"TotalTickets": event.TotalTickets,
			"SoldTickets": soldTickets,
		},
		"IsAuthenticated": isAuthenticated,
		"Username": username,
		"HasTicket": hasTicket,
	}

	log.Printf("Rendering template with data: %+v", data)

	c.HTML(http.StatusOK, "event.html", data)
}

func GetEvents(c *gin.Context) {
	var events []models.Event
	if err := db.DB.Preload("Category").Preload("Organizer").Find(&events).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при загрузке событий",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"events": events,
	})
}

func BuyTicket(c *gin.Context) {
	eventID := c.Param("id")
	userID := c.GetUint("user_id")

	var event models.Event
	if err := db.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": "Событие не найдено",
		})
		return
	}

	// Проверяем, не превышено ли количество билетов
	if event.SoldTickets >= event.TotalTickets {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": "Все билеты проданы",
		})
		return
	}

	// Проверяем, не купил ли пользователь уже билет на это мероприятие
	var existingTicket models.UserTicket
	if err := db.DB.Where("user_id = ? AND event_id = ?", userID, event.ID).First(&existingTicket).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": "Вы уже купили билет на это мероприятие",
		})
		return
	}

	// Создаем запись в Ticket
	ticket := models.Ticket{
		EventID: event.ID,
		UserID:  userID,
	}

	// Создаем запись в UserTicket
	userTicket := models.UserTicket{
		UserID:    userID,
		EventID:   event.ID,
		Quantity:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Начинаем транзакцию
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": "Ошибка при создании билета",
		})
		return
	}

	// Создаем запись в Ticket
	if err := tx.Create(&ticket).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": "Ошибка при создании билета",
		})
		return
	}

	// Создаем запись в UserTicket
	if err := tx.Create(&userTicket).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": "Ошибка при создании записи о билете",
		})
		return
	}

	// Увеличиваем счетчик проданных билетов
	if err := tx.Model(&event).Update("sold_tickets", event.SoldTickets+1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": "Ошибка при обновлении количества билетов",
		})
		return
	}

	// Завершаем транзакцию
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": "Ошибка при сохранении транзакции",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Билет успешно куплен!",
	})
}
