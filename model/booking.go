package model

import (
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID     uint         `gorm:"type:bigint unsigned;not null" json:"user_id"`
	User       User         `gorm:"foreignKey:UserID;references:ID" json:"user"`
	TicketID   uint         `gorm:"type:bigint unsigned;not null" json:"ticket_id"`
	Ticket     TicketMaster `gorm:"foreignKey:TicketID;references:ID" json:"ticket"`
	Quantity   int          `gorm:"not null" json:"quantity"`
	TotalPrice int          `gorm:"not null" json:"total_price"`
}
