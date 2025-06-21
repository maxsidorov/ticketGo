package service

import (
	"errors"
	"github.com/maxsidorov/ticketGo/models"
	"github.com/maxsidorov/ticketGo/storage"
	"gorm.io/gorm"
	"time"
)

type EventService struct {
	storage *storage.EventStorage
	db      *gorm.DB
}

func (s *EventService) CreateEvent(event *models.Event) (uint, error) {
	if event.DateTime.Before(time.Now()) {
		return 0, errors.New("event date must be in the future")
	}
	return s.storage.Create(event)
}

func (s *EventService) GetEvents(page, pageSize int) ([]models.Event, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.storage.GetAll(page, pageSize)
}

func (s *EventService) GetEvent(id int) (*models.Event, error) {
	return s.storage.GetByID(id)
}

func (s *EventService) UpdateEvent(event *models.Event) error {
	if event.DateTime.Before(time.Now()) {
		return errors.New("cannot update past events")
	}
	return s.storage.Update(event)
}

func (s *EventService) DeleteEvent(id int) error {
	return s.storage.Delete(id)
}

func (s *EventService) GetUserTicketsCount(eventID, userID int) (int, error) {
	var count int64
	err := s.db.Model(&models.UserTicket{}).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
