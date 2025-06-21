package storage

import (
	"github.com/maxsidorov/ticketGo/models"
	"gorm.io/gorm"
)

type EventStorage struct {
	db *gorm.DB
}

func (s *EventStorage) Create(event *models.Event) (uint, error) {
	if err := s.db.Create(event).Error; err != nil {
		return 0, err
	}
	return event.ID, nil
}

func (s *EventStorage) GetAll(page, pageSize int) ([]models.Event, error) {
	var events []models.Event
	err := s.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&events).Error
	return events, err
}

func (s *EventStorage) GetByID(id int) (*models.Event, error) {
	var event models.Event
	if err := s.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *EventStorage) Update(event *models.Event) error {
	return s.db.Save(event).Error
}

func (s *EventStorage) Delete(id int) error {
	return s.db.Delete(&models.Event{}, id).Error
}
