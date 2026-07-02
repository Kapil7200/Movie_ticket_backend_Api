package model

import (
	"gorm.io/gorm"
)

type TicketMaster struct {
	gorm.Model
	MovieName       string  `gorm:"not null" json:"movie_name"`
	TotalTicket     int     `gorm:"not null" json:"total_ticket"`
	AvailableTicket int     `gorm:"not null" json:"available_ticket"`
	PricePerTicket  float64 `gorm:"not null" json:"price_per_ticket"`
}
