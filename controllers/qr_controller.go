package controllers

import (
	"encoding/json"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/models"
	"image/png"
	"net/http"
	"strconv"
)

func GenerateTicketQR(c *gin.Context) {
	// Получаем ID билета из URL
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Получаем данные билета из базы
	var ticket models.UserTicket
	var user models.User
	var event models.Event
	if err := db.DB.Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	if err := db.DB.Where("id = ?", ticket.UserID).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := db.DB.Where("id = ?", ticket.EventID).First(&event).Error; err != nil {
		return
	}
	// Создаем данные для QR-кода (можно настроить по вашему усмотрению)
	qrData := map[string]interface{}{
		"ticket_id":   ticket.ID,
		"event_id":    ticket.EventID,
		"user_id":     ticket.UserID,
		"valid_until": event.DateTime, // или другая дата\
		"username":    user.Username,
		"email":       user.Email,
	}

	// Конвертируем данные в JSON
	jsonData, err := json.Marshal(qrData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR data"})
		return
	}

	// Генерируем QR-код
	qrCode, err := qr.Encode(string(jsonData), qr.M, qr.Auto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Масштабируем QR-код (по желанию)
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scale QR code"})
		return
	}

	// Устанавливаем заголовки для скачивания
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "attachment; filename=ticket_qr.png")

	// Отправляем PNG изображение
	if err := png.Encode(c.Writer, qrCode); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode QR code"})
		return
	}
}
