package repository

import (
	"errors"

	"movie_ticket/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GORMRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) *GORMRepository {
	return &GORMRepository{db: db}
}

func (r *GORMRepository) CreateUser(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		logrus.Infof("CreateUser failed. UserName=%s Email=%s Error=%s", user.UserName, user.Email, err)
		return err
	}

	return nil
}

func (r *GORMRepository) GetUserByUserName(userName string) (*model.User, error) {
	user := &model.User{}

	if err := r.db.Where("user_name = ?", userName).First(user).Error; err != nil {
		logrus.Infof("GetUserByUserName failed. UserName=%s Error=%v", userName, err)
		return nil, err
	}

	return user, nil
}

func (r *GORMRepository) GetUserByID(id uint) (*model.User, error) {
	user := &model.User{}
	if err := r.db.First(user, id).Error; err != nil {
		logrus.Infof("GetUserByID failed. UserID=%d Error=%v", id, err)
		return nil, err

	}
	return user, nil
}

func (r *GORMRepository) ListUsers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		logrus.Infof("ListUsers failed. Error=%v", err)
		return nil, err
	}
	return users, nil
}

func (r *GORMRepository) CreateMoviesTicket(ticket *model.TicketMaster) error {

	if err := r.db.Create(ticket).Error; err != nil {

		logrus.Infof("CreateMoviesTicket failed. TicketID=%d Movie=%s Error=%v",
			ticket.ID, ticket.MovieName, err)

		return err
	}

	logrus.Infof("CreateMoviesTicket successful. TicketID=%d Movie=%s",
		ticket.ID, ticket.MovieName)

	return nil
}

func (r *GORMRepository) GetTicketByID(id uint) (*model.TicketMaster, error) {
	ticket := &model.TicketMaster{}
	if err := r.db.First(ticket, id).Error; err != nil {
		logrus.Infof("GetTicketByID failed. TicketID=%d Error=%v", id, err)
		return nil, err
	}
	return ticket, nil
}

func (r *GORMRepository) GetAllTickets() ([]model.TicketMaster, error) {
	var tickets []model.TicketMaster
	if err := r.db.Find(&tickets).Error; err != nil {
		logrus.Infof("GetAllTickets failed. Error=%v", err)
		return nil, err
	}
	return tickets, nil
}

func (r *GORMRepository) UpdateTicket(ticket *model.TicketMaster) error {
	result := r.db.Model(&model.TicketMaster{}).Where("id = ?", ticket.ID).Updates(map[string]interface{}{
		"movie_name":       ticket.MovieName,
		"total_ticket":     ticket.TotalTicket,
		"available_ticket": ticket.AvailableTicket,
		"price_per_ticket": ticket.PricePerTicket,
	})
	if result.Error != nil {
		logrus.Infof("UpdateTicket failed. TicketID=%d Error=%v", ticket.ID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		logrus.Infof("UpdateTicket failed. TicketID=%d Reason=ticket not found", ticket.ID)
		return errors.New("ticket not found")
	}
	return nil
}

func (r *GORMRepository) BookTicket(ticketID, userID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		logrus.Infof("BookTicket: transaction start failed: %v", tx.Error)
		return tx.Error
	}

	ticket := &model.TicketMaster{}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(ticket, ticketID).Error; err != nil {
		logrus.Infof("BookTicket: ticket %d not found: %v", ticketID, err)
		tx.Rollback()
		return err
	}

	if ticket.AvailableTicket < quantity {
		logrus.Infof(
			"BookTicket: insufficient tickets. TicketID=%d Available=%d Requested=%d", ticketID, ticket.AvailableTicket, quantity,
		)
		tx.Rollback()
		return errors.New("not enough tickets available")
	}

	ticket.AvailableTicket -= quantity
	if err := tx.Save(ticket).Error; err != nil {
		logrus.Infof("BookTicket: failed to update ticket %d: %v", ticketID, err)
		tx.Rollback()
		return err
	}
	totalPrice := int(ticket.PricePerTicket) * (quantity)

	booking := &model.Booking{
		UserID:     userID,
		TicketID:   ticketID,
		Quantity:   quantity,
		TotalPrice: totalPrice,
	}
	if err := tx.Create(booking).Error; err != nil {
		logrus.Infof("BookTicket failed. UserID=%d TicketID=%d Error=%v", userID, ticketID, err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		logrus.Infof("BookTicket: commit failed: %v", err)
		return err
	}

	return nil
}

func (r *GORMRepository) GetBookingsByUserID(userID uint) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.db.Preload("Ticket").Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *GORMRepository) GetUserBookings(userID uint) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.db.Preload("Ticket").Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		logrus.Infof("GetUserBookings failed. UserID=%d Error=%v", userID, err)
	}
	return bookings, nil
}
