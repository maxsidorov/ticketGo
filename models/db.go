package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Автоматическая миграция моделей
	err = DB.AutoMigrate(&User{}, &Event{}, &Ticket{})
	if err != nil {
		return err
	}

	return nil
} 