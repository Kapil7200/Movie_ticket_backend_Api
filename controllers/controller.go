package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"movie_ticket/dto"
	"movie_ticket/middleware"
	"movie_ticket/service"
	"movie_ticket/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	svc     service.ServiceInterface
	jwtUtil *middleware.JWTUtil
}

func NewController(svc service.ServiceInterface, jwtUtil *middleware.JWTUtil) *Controller {
	return &Controller{
		svc:     svc,
		jwtUtil: jwtUtil,
	}
}

func (c *Controller) RegisterUser(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("RegisterUser@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	user, err := c.svc.RegisterUser(req.UserName, req.Password, req.Email)
	if err != nil {
		logrus.Info("Register user failed: %v", err)
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

func (c *Controller) UserLogin(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("UserLogin@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Info("Login failed. Error=%v", err)
		utils.ValidationError(ctx, err)
		return
	}

	user, err := c.svc.Authenticate(req.UserName, req.Password)
	if err != nil {
		logrus.Info("Login failed for user %s: %v", req.UserName, err)
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
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("VerifyToken@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
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

func (c *Controller) ListUsers(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("ListUsers@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
	users, err := c.svc.ListUsers()
	if err != nil {
		logrus.Info("List users failed: %v", err)
		utils.ValidationError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, users)
}
func (c *Controller) CreateMoviesTicket(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("CreateMoviesTicket@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
	var req dto.TicketRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Info("CreateMoviesTicket failed. Error=%v", err)
		utils.ValidationError(ctx, err)
		return
	}

	ticket, err := c.svc.CreateMoviesTicket(req)
	if err != nil {
		logrus.Info("CreateMoviesTicket failed. Error=%v", err)
		utils.ValidationError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, ticket)
}

func (c *Controller) GetAllTickets(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("GetAllTickets@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
	tickets, err := c.svc.GetAllTickets()
	if err != nil {
		logrus.Info("Get all tickets failed: %v", err)
		utils.ValidationError(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, tickets)
}

func (c *Controller) GetTicketByID(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("GetTicketByID@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.ValidationError(ctx, err)
		return
	}

	ticket, err := c.svc.GetTicketByID(uint(id))
	if err != nil {
		logrus.Info("Get ticket by ID %d failed: %v", id, err)
		utils.ValidationError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, ticket)
}

func (c *Controller) UpdateTicket(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("UpdateTicket@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
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

	if err := c.svc.UpdateTicket(uint(id), req); err != nil {
		logrus.Info("Update ticket %d failed: %v", id, err)
		utils.ValidationError(ctx, err)
		return
	}

	updated, err := c.svc.GetTicketByID(uint(id))
	if err != nil {
		logrus.Info("Get updated ticket %d failed: %v", id, err)
		utils.ValidationError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, updated)
}

func (c *Controller) BookTicket(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("BookTicket@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
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
		logrus.Info("Booking failed. UserID=%d TicketID=%d Error=%v", claims.UserID, id, err)
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
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("GetBookings@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()

	claims := middleware.GetClaims(ctx)
	if claims == nil {
		utils.ValidationError(ctx, fmt.Errorf("authorization required"))
		return
	}

	bookings, err := c.svc.GetBookingsByUserID(claims.UserID)
	if err != nil {
		logrus.Infof("Get bookings failed for user %d: %v", claims.UserID, err)
		utils.ValidationError(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, bookings)
}

func (c *Controller) GetUserByID(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info("GetUserByID@panicInfo:", r)
			utils.InternalServerErrorResponse(ctx, fmt.Errorf("%v", r))
		}
	}()
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
