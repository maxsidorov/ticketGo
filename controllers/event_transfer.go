package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/models"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ExportEvent(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	var events []models.Event
	if err := db.DB.Where("admin_id = ?", userID).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении событий"})
		return
	}

	// Создаем временный файл
	fileName := fmt.Sprintf("events_export_%d.json", time.Now().Unix())
	filePath := filepath.Join(os.TempDir(), fileName)

	file, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании файла"})
		return
	}
	defer file.Close()

	// Пишем данные в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
		return
	}

	// Отправляем файл пользователю
	c.FileAttachment(filePath, fileName)

	// Удаляем временный файл после отправки
	defer os.Remove(filePath)
}
