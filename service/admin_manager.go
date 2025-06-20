package service

import (
	"github.com/maxsidorov/ticketGo/db"
	"github.com/maxsidorov/ticketGo/models"
	"log"
)

func AddEvent(event models.Event) {
	if err := db.DB.Create(&event).Error; err != nil {
		log.Println(err)
	}
}

func UpdateEvent(event models.Event) error {
	return db.DB.Save(&event).Error
}
