package models

type Ticket struct {
	ID      uint `gorm:"primaryKey"`
	EventID uint `gorm:"not null"`
	Event   Event `gorm:"foreignKey:EventID"`
} 