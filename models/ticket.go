package models

type Ticket struct {
	ID      uint   `gorm:"primaryKey"`
	EventID uint   `gorm:"not null"`
	Event   Event  `gorm:"foreignKey:EventID"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID"`
	QR      string `gorm:"not null"`
}
