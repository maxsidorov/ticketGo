package storage

import (
  "database/sql"
  "fmt"
  "time"
)
// инициализация базы данных
func InitDB(file string) (*sql.DB, error) {
  db, err := sql.Open("sqlite3", file)
  if err != nil {
    return nil, fmt.Errorf("failed to open database: %w", err)
  }

  // Применяем схему
  if _, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS events (...); -- из schema.sql
  `); err != nil {
    return nil, fmt.Errorf("failed to create tables: %w", err)
  }

  return db, nil
}
// создание события
type EventStorage struct {
  db *sql.DB
}
func NewEventStorage(db *sql.DB) *EventStorage {
  return &EventStorage{db: db}
}
func (s *EventStorage) Create(event *Event) (int, error) {
  // Реализация простого создания
  res, err := s.db.Exec(
    `INSERT INTO events (id, title, date, place, decsription, price, tickets, sold_tickets, image, discount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
    event.ID, event.Title, event.Date, event.Place, event.Decsription, event.Price, event.Tickets, event.Sold_tickets, event.Image, event.Discount,
  )
  if err != nil {
    return 0, err
  }

  id, _ := res.LastInsertId()
  return int(id), nil
}

  func (s *EventStorage) GetAll() ([]Event, error) {
    // Базовая реализация без пагинации
    rows, err := s.db.Query("SELECT id, title, date, place, decsription, price, tickets, sold_tickets, image, discount FROM events")
    if err != nil {
      return nil, err
    }
    defer rows.Close()
    var events []Event
    for rows.Next() {
      var e Event
      if err := rows.Scan(&e.ID, &e.Title, &e.Date, &e.Place, &e.Decsription, &e.Price, &e.Tickets, &e.Sold_tickets, &e.Image, &e.Discount); err != nil {
        return nil, err
      }
      events = append(events, e)
    }
    return events, nil
  }
// Модель данных
type Event struct {
  ID        int
  Title     string
  Date      time.Time
  Place  string
  Decsription string
  Price     float64
  Tickets   int
  Sold_tickets int
  Image     string
  Discount  float64
}
