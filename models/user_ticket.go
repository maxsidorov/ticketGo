package models

type UserTicket struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	EventID uint `gorm:"not null"`
	User   User  `gorm:"foreignKey:UserID"`
	Event  Event `gorm:"foreignKey:EventID"`
} 