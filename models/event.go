package models

import (
	"time"
)

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
