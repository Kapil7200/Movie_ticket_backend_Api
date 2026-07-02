package routes

import (
	"movie_ticket/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, controller *controllers.Controller, authMiddleware gin.HandlerFunc) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	authGroup := router.Group("/")
	authGroup.Use(authMiddleware)
	{
		authGroup.GET("/bookings", controller.GetBookings)
		authGroup.GET("/users/:id", controller.GetUserByID)
		authGroup.GET("/tickets", controller.GetAllTickets)
		authGroup.POST("/create-tickets", controller.CreateTicket)
		authGroup.GET("/tickets/:id", controller.GetTicketByID)
		authGroup.PUT("/update-tickets/:id", controller.UpdateTicket)
		authGroup.POST("/book-tickets/:id/book", controller.BookTicket)
	}
}
