package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"movie_ticket/dto"
	"movie_ticket/middleware"
	"movie_ticket/model"
	"movie_ticket/service"
	"movie_ticket/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	svc     *service.Service
	jwtUtil *middleware.JWTUtil
}

func NewController(svc *service.Service, jwtUtil *middleware.JWTUtil) *Controller {
	return &Controller{svc: svc, jwtUtil: jwtUtil}
}

func (c *Controller) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	user, err := c.svc.RegisterUser(req.UserName, req.Password, req.Email)
	if err != nil {
		utils.Logger.Printf("Register user failed: %v", err)
		utils.ValidationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"user_name":  user.UserName,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (c *Controller) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	user, err := c.svc.Authenticate(req.UserName, req.Password)
	if err != nil {
		utils.Logger.Printf("Login failed for user %s: %v", req.UserName, err)
		utils.ValidationError(ctx, err)
		return
	}

	token, err := c.jwtUtil.CreateToken(user.ID, user.UserName)
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *Controller) VerifyToken(ctx *gin.Context) {
	claims := middleware.GetClaims(ctx)
	if claims == nil {
		utils.ValidationError(ctx, fmt.Errorf("authorization required"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":   claims.UserID,
		"user_name": claims.UserName,
	})
}

func (c *Controller) CreateTicket(ctx *gin.Context) {
	var req dto.TicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Logger.Printf("CreateTicket failed. Error=%v", err)
		utils.ValidationError(ctx, err)
		return
	}

	ticket := &model.TicketMaster{
		MovieName:       req.MovieName,
		TotalTicket:     req.TotalTicket,
		AvailableTicket: req.TotalTicket,
		PricePerTicket:  req.PricePerTicket,
	}

	if err := c.svc.CreateTicket(ticket); err != nil {
		utils.Logger.Printf("Create ticket failed: %v", err)
		utils.ValidationError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, ticket)
}

func (c *Controller) GetAllTickets(ctx *gin.Context) {
	tickets, err := c.svc.GetAllTickets()
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, tickets)
}

func (c *Controller) GetTicketByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	ticket, err := c.svc.GetTicketByID(uint(id))
	if err != nil {
		utils.Logger.Printf("Get ticket by ID %d failed: %v", id, err)
		utils.ValidationError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, ticket)
}

func (c *Controller) UpdateTicket(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	var req dto.TicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	ticket := &model.TicketMaster{
		Model:           gorm.Model{ID: uint(id)},
		MovieName:       req.MovieName,
		TotalTicket:     req.TotalTicket,
		AvailableTicket: req.TotalTicket,
		PricePerTicket:  req.PricePerTicket,
	}

	if err := c.svc.UpdateTicket(ticket); err != nil {
		utils.Logger.Printf("Update ticket %d failed: %v", id, err)
		utils.ValidationError(ctx, err)
		return
	}

	utils.Logger.Printf("Ticket updated successfully: %d", id)
	updated, err := c.svc.GetTicketByID(uint(id))
	if err != nil {
		utils.Logger.Printf("Get updated ticket by ID %d failed: %v", id, err)
		utils.ValidationError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, updated)
}

func (c *Controller) BookTicket(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	var req dto.BookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	claims := middleware.GetClaims(ctx)
	if claims == nil {
		utils.ValidationError(ctx, fmt.Errorf("authorization required"))
		return
	}

	if err := c.svc.BookTicket(uint(id), claims.UserID, req.Quantity); err != nil {
		utils.Logger.Printf("Booking failed. UserID=%d TicketID=%d Error=%v", claims.UserID, id, err)
		utils.ValidationError(ctx, err)
		return
	}

	bookings, err := c.svc.GetBookingsByUserID(claims.UserID)
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	var response []dto.BookingResponse

	for _, booking := range bookings {
		response = append(response, dto.BookingResponse{
			UserID:     booking.UserID,
			TicketID:   booking.TicketID,
			Quantity:   booking.Quantity,
			TotalPrice: booking.TotalPrice,
		})
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"message":  "booking successful",
		"bookings": response,
	})
}

func (c *Controller) GetBookings(ctx *gin.Context) {
	claims := middleware.GetClaims(ctx)
	if claims == nil {
		utils.ValidationError(ctx, fmt.Errorf("authorization required"))
		return
	}

	bookings, err := c.svc.GetBookingsByUserID(claims.UserID)
	if err != nil {
		utils.Logger.Printf("Get bookings failed for user %d: %v", claims.UserID, err)
		utils.ValidationError(ctx, err)
		return
	}
	var response []dto.BookingResponse

	for _, booking := range bookings {
		response = append(response, dto.BookingResponse{
			UserID:     booking.UserID,
			TicketID:   booking.TicketID,
			Quantity:   booking.Quantity,
			TotalPrice: booking.TotalPrice,
		})
	}

	utils.SuccessResponse(ctx, http.StatusOK, response)
}

func (c *Controller) GetUserByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	user, err := c.svc.GetUserByID(uint(id))
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	bookings, err := c.svc.GetBookingsByUserID(uint(id))
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{"user": user, "bookings": bookings})
}
