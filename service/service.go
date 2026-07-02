package service

import (
	"movie_ticket/model"
	"movie_ticket/repository"
	"movie_ticket/utils"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(userName, password, email string) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		utils.Logger.Printf("Password hashing failed: %v", err)
		return nil, err
	}

	user := &model.User{
		UserName: userName,
		Password: hashedPassword,
		Email:    email,
	}

	if err := s.repo.CreateUser(user); err != nil {
		utils.Logger.Printf(" failed to Create user: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *Service) Authenticate(userName, password string) (*model.User, error) {
	user, err := s.repo.GetUserByUserName(userName)
	if err != nil {
		utils.Logger.Printf("User not found: %s", userName)
		return nil, err
	}

	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		utils.Logger.Printf("Invalid password for user: %s", userName)
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateTicket(ticket *model.TicketMaster) error {
	if ticket.AvailableTicket == 0 {
		ticket.AvailableTicket = ticket.TotalTicket
	}

	if err := s.repo.CreateTicket(ticket); err != nil {
		utils.Logger.Printf("Create ticket failed: %v", err)
		return err
	}

	return nil
}

func (s *Service) GetAllTickets() ([]model.TicketMaster, error) {
	tickets, err := s.repo.GetAllTickets()
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (s *Service) GetTicketByID(id uint) (*model.TicketMaster, error) {
	ticket, err := s.repo.GetTicketByID(id)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *Service) UpdateTicket(ticket *model.TicketMaster) error {
	if err := s.repo.UpdateTicket(ticket); err != nil {
		return err
	}

	return nil
}

func (s *Service) BookTicket(ticketID, userID uint, quantity int) error {
	if err := s.repo.BookTicket(ticketID, userID, quantity); err != nil {
		utils.Logger.Printf("Book ticket failed. Ticket=%d User=%d Error=%v",ticketID, userID, err)
		return err
	}

	return nil
}

func (s *Service) GetUserByID(id uint) (*model.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetBookingsByUserID(userID uint) ([]model.Booking, error) {
	bookings, err := s.repo.GetBookingsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}
