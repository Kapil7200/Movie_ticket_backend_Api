package routes

import (
	"movie_ticket/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, controller *controllers.Controller, authMiddleware gin.HandlerFunc) {
	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.UserLogin)

	authGroup := router.Group("/")
	authGroup.Use(authMiddleware)
	{
		authGroup.GET("/bookings-listing", controller.GetBookings)
		//authGroup.GET("/users-list/:id", controller.GetUserByID)
		authGroup.GET("/users-listing", controller.ListUsers)
		authGroup.GET("/tickets-listing", controller.GetAllTickets)
		authGroup.POST("/create-tickets", controller.CreateMoviesTicket)
		authGroup.GET("/tickets-list/:id", controller.GetTicketByID)
		authGroup.PUT("/update/ticket/:id", controller.UpdateTicket)
		authGroup.POST("/book-ticket/:id/book", controller.BookTicket)
	}
}
