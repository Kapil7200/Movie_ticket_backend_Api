package repository

import (
	"movie_ticket/dto"
	"movie_ticket/model"
)

type Repository interface {
	CreateUser(user *model.User) error
	GetUserByUserName(userName string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	CreateMoviesTicket(ticket *model.TicketMaster) error
	GetTicketByID(id uint) (*model.TicketMaster, error)
	GetAllTickets() ([]model.TicketMaster, error)
	UpdateTicket(ticket *model.TicketMaster) error
	BookTicket(ticketID, userID uint, quantity int) error
	GetBookingsByUserID(userID uint) ([]dto.BookingResponse, error)
	GetUserBookings(userID uint) ([]model.Booking, error)
	ListUsers() ([]model.User, error)

}
