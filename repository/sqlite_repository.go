package repository

import (
	"errors"

	"movie_ticket/model"
	"movie_ticket/utils"

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
		utils.Logger.Printf("CreateUser failed. UserName=%s Email=%s Error=%v", user.UserName, user.Email, err)
		return err
	}

	return nil
}

func (r *GORMRepository) GetUserByUserName(userName string) (*model.User, error) {
	user := &model.User{}

	if err := r.db.Where("user_name = ?", userName).First(user).Error; err != nil {
		utils.Logger.Printf("GetUserByUserName failed. UserName=%s Error=%v", userName, err)
		return nil, err
	}

	return user, nil
}

func (r *GORMRepository) GetUserByID(id uint) (*model.User, error) {
	user := &model.User{}
	if err := r.db.First(user, id).Error; err != nil {
		utils.Logger.Printf("GetUserByID failed. UserID=%d Error=%v", id, err)
		return nil, err

	}
	return user, nil
}

func (r *GORMRepository) CreateTicket(ticket *model.TicketMaster) error {
	if ticket.AvailableTicket == 0 {
		ticket.AvailableTicket = ticket.TotalTicket
	}
	if err := r.db.Create(ticket).Error; err != nil {
		logrus.Error("CreateTicket failed. TicketID=%d  Movie=%s Error=%v", ticket.ID, ticket.MovieName, err)
		utils.Logger.Printf("CreateTicket failed. TicketID=%d  Movie=%s Error=%v", ticket.ID, ticket.MovieName, err)
		return err
	}
	return nil
}

func (r *GORMRepository) GetTicketByID(id uint) (*model.TicketMaster, error) {
	ticket := &model.TicketMaster{}
	if err := r.db.First(ticket, id).Error; err != nil {
		utils.Logger.Printf("GetTicketByID failed. TicketID=%d Error=%v", id, err)
		return nil, err
	}
	return ticket, nil
}

func (r *GORMRepository) GetAllTickets() ([]model.TicketMaster, error) {
	var tickets []model.TicketMaster
	if err := r.db.Find(&tickets).Error; err != nil {
		utils.Logger.Printf("GetAllTickets failed. Error=%v", err)
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
		logrus.Error("UpdateTicket failed. TicketID=%d Error=%v", ticket.ID, result.Error)
		utils.Logger.Printf("UpdateTicket failed. TicketID=%d Error=%v", ticket.ID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		utils.Logger.Printf("UpdateTicket failed. TicketID=%d Reason=ticket not found", ticket.ID)
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
		utils.Logger.Printf("BookTicket: transaction start failed: %v", tx.Error)
		return tx.Error
	}

	ticket := &model.TicketMaster{}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(ticket, ticketID).Error; err != nil {
		utils.Logger.Printf("BookTicket: ticket %d not found: %v", ticketID, err)
		tx.Rollback()
		return err
	}

	if ticket.AvailableTicket < quantity {
		utils.Logger.Printf(
			"BookTicket: insufficient tickets. TicketID=%d Available=%d Requested=%d", ticketID, ticket.AvailableTicket, quantity,
		)
		tx.Rollback()
		return errors.New("not enough tickets available")
	}

	ticket.AvailableTicket -= quantity
	if err := tx.Save(ticket).Error; err != nil {
		utils.Logger.Printf("BookTicket: failed to update ticket %d: %v", ticketID, err)
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
		utils.Logger.Printf("BookTicket failed. UserID=%d TicketID=%d Error=%v", userID, ticketID, err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		utils.Logger.Printf("BookTicket: commit failed: %v", err)
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
		return nil, err
	}
	return bookings, nil
}
