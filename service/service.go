package service

import (
	"movie_ticket/dto"
	"movie_ticket/model"
	"movie_ticket/repository"
	"movie_ticket/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
		logrus.Info("Password hashing failed: %v", err)
		return nil, err
	}

	user := &model.User{
		UserName: userName,
		Password: hashedPassword,
		Email:    email,
	}

	if err := s.repo.CreateUser(user); err != nil {
		logrus.Info(" failed to Create user: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *Service) Authenticate(userName, password string) (*model.User, error) {
	user, err := s.repo.GetUserByUserName(userName)
	if err != nil {
		logrus.Info("User not found: %s", userName)
		return nil, err
	}

	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		logrus.Info("Invalid password for user: %s", userName)
		return nil, err
	}

	return user, nil
}

func (s *Service) ListUsers() ([]model.User, error) {
	users, err := s.repo.ListUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) CreateMoviesTicket(req dto.TicketRequest) (*model.TicketMaster, error) {

	ticket := &model.TicketMaster{
		MovieName:       req.MovieName,
		TotalTicket:     req.TotalTicket,
		AvailableTicket: req.TotalTicket,
		PricePerTicket:  req.PricePerTicket,
	}

	if err := s.repo.CreateMoviesTicket(ticket); err != nil {
		logrus.Info("CreateMoviesTicket failed. Error=%v", err)
		return nil, err
	}

	return ticket, nil
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
func (s *Service) UpdateTicket(id uint, req dto.TicketRequest) error {

	ticket := &model.TicketMaster{
		Model: gorm.Model{
			ID: id,
		},
		MovieName:       req.MovieName,
		TotalTicket:     req.TotalTicket,
		AvailableTicket: req.TotalTicket,
		PricePerTicket:  req.PricePerTicket,
	}

	return s.repo.UpdateTicket(ticket)
}

func (s *Service) BookTicket(ticketID, userID uint, quantity int) error {
	if err := s.repo.BookTicket(ticketID, userID, quantity); err != nil {
		logrus.Info("Book ticket failed. Ticket=%d User=%d Error=%v", ticketID, userID, err)
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
