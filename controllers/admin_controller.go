package controllers

import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/service"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func AdminPage(c *gin.Context) {
	var events []models.Event
	session := sessions.Default(c)
	userID := session.Get("user_id")
	query := db.DB.Model(&models.Event{})
	query = query.Where("admin_id = ?", userID)

	username := session.Get("username")
	var user models.User
	db.DB.Where("username = ?", username).First(&user)
	if user.AdminLevel == 0 {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"IsAuthenticated": c.GetBool("is_authenticated"),
		})
		return
	}

	search := c.Query("search")
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Find(&events)

	if eventID := c.Query("delete"); eventID != "" {
		id, _ := strconv.Atoi(eventID)
		db.DB.Where("event_id = ?", id).Delete(&models.UserTicket{})
		db.DB.Delete(&models.Event{}, id)
	}
	query.Find(&events)
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"events": events,
	})
}

func AddEventPage(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Обработка GET запроса (открытие формы)
		eventID := c.Query("edit")
		var event models.Event

		if eventID != "" {
			// Режим редактирования - загружаем существующее событие
			id, err := strconv.Atoi(eventID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID события"})
				return
			}

			if err := db.DB.First(&event, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Событие не найдено"})
				return
			}
		}
		// Передаем событие в шаблон (пустое или заполненное)
		c.HTML(http.StatusOK, "event_add.html", gin.H{
			"event":  event,
			"isEdit": eventID != "",
		})
		return
	}
	if c.Request.Method == "POST" {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		// Получаем данные из формы
		title := c.PostForm("title")
		datetimeStr := c.PostForm("datetime")
		location := c.PostForm("location")
		description := c.PostForm("description")
		priceStr := c.PostForm("price")
		totalTicketsStr := c.PostForm("total_tickets")
		category := c.PostForm("category")
		eventID := c.PostForm("event_id")

		// Парсим числовые значения
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная цена"})
			return
		}

		totalTickets, err := strconv.Atoi(totalTicketsStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректное количество билетов"})
			return
		}

		// Парсим дату и время
		datetime, err := time.Parse("2006-01-02T15:04", datetimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат даты и времени"})
			return
		}
		datetime = time.Date(
			datetime.Year(),
			datetime.Month(),
			datetime.Day(),
			datetime.Hour(),
			datetime.Minute(),
			datetime.Second(),
			datetime.Nanosecond(),
			time.Local,
		)

		fileHeader, err := c.FormFile("image")
		var uploadPath string
		if err == nil {
			// Открываем загруженный файл
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при открытии файла"})
				return
			}
			defer file.Close()

			// Декодируем изображение (поддержка PNG, JPEG, GIF)
			img, _, err := image.Decode(file)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неподдерживаемый формат изображения"})
				return
			}

			// Генерируем уникальное имя файла
			hash := sha1.New()
			hash.Write([]byte(title + time.Now().String()))
			filename := fmt.Sprintf("%x.jpeg", hash.Sum(nil))

			// Создаем выходной файл JPEG
			uploadPath = filepath.Join("static", filename)
			outFile, err := os.Create(uploadPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании файла"})
				return
			}
			defer outFile.Close()

			// Кодируем в JPEG с качеством 80%
			jpegOptions := jpeg.Options{Quality: 80}
			if err := jpeg.Encode(outFile, img, &jpegOptions); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при конвертации изображения"})
				return
			}
		}
		// Обработка загрузки изображения

		// Создаем структуру Event
		newEvent := models.Event{
			Title:        title,
			DateTime:     datetime,
			Location:     location,
			Description:  description,
			Price:        price,
			TotalTickets: totalTickets,
			SoldTickets:  0,
			Category:     category,
			AdminID:      userID.(uint),
			Image:        "/" + filepath.ToSlash(uploadPath),
		}
		fmt.Println(eventID, "hello")
		if eventID != "" {
			// Режим редактирования - обновляем существующее
			id, _ := strconv.Atoi(eventID)
			newEvent.ID = uint(id)
			service.UpdateEvent(newEvent)
		} else {
			// Режим создания - добавляем новое
			service.AddEvent(newEvent)
		}
		// Перенаправляем после успешного сохранения
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	c.HTML(http.StatusOK, "event_add.html", gin.H{})
}

func AdminUsersPage(c *gin.Context) {
	searchQuery := c.Query("search")
	var users []models.User

	query := db.DB.Model(&models.User{})

	session := sessions.Default(c)
	username := session.Get("username")
	var user models.User
	db.DB.Where("username = ?", username).First(&user)
	if user.AdminLevel != 2 {
		var events []models.Event
		queryDop := db.DB.Model(&models.Event{})
		queryDop = queryDop.Where("admin_id = ?", user.ID)
		queryDop.Find(&events)
		c.HTML(http.StatusBadRequest, "admin.html", gin.H{
			"events": events,
		})
		return
	}

	if searchQuery != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%")
	}
	query.Find(&users)

	c.HTML(http.StatusOK, "admin_users.html", gin.H{
		"Users":       users,
		"SearchQuery": searchQuery,
		"Success":     c.Query("success"),
		"Error":       c.Query("error"),
	})
}

func UpdateUserAdminLevel(c *gin.Context) {
	userID := c.PostForm("user_id")
	adminLevel, err := strconv.Atoi(c.PostForm("admin_level"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/users?error=Неверный уровень доступа")
		return
	}

	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("admin_level", adminLevel).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/users?error=Ошибка обновления пользователя")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users?success=Уровень доступа обновлен")
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := db.DB.Delete(&models.User{}, userID).Error; err != nil {
		c.Redirect(http.StatusFound, "/admin/users?error=Ошибка удаления пользователя")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users?success=Пользователь удален")
}
