package models

import (
	"time"
)

type Event struct {
	ID           uint      `gorm:"primaryKey"`
	Title        string    `gorm:"not null"`
	Image        string    `gorm:"not null"`
	DateTime     time.Time `gorm:"column:date_time;not null"`
	Location     string    `gorm:"not null"`
	Description  string    `gorm:"not null"`
	Price        float64   `gorm:"not null"`
	TotalTickets int       `gorm:"not null"`
	SoldTickets  int       `gorm:"default:0"`
	Category     string    `gorm:"not null;default:'concert'"`
}
