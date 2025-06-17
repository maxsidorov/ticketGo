package models

import (
	"time"
)

type User struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Username   string         `gorm:"unique;not null" json:"username"`
	Email      string         `gorm:"unique;not null" json:"email"`
	Password   string         `gorm:"not null" json:"-"`
	AdminLevel int            `gorm:"default:0" json:"admin_level"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Tickets    []UserTicket   `gorm:"foreignKey:UserID" json:"tickets"`
}

type UserTicket struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	EventID      uint      `gorm:"not null" json:"event_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	Event        Event     `gorm:"foreignKey:EventID" json:"event"`
	TicketsCount int       `gorm:"not null;default:1" json:"tickets_count"`
	Quantity     int       `gorm:"not null;default:1" json:"quantity"`
	Status       string    `gorm:"not null;default:'active'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
} 