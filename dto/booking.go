package dto

type BookingRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}
type BookingResponse struct {
	UserID     uint `json:"user_id"`
	TicketID   uint `json:"ticket_id"`
	Quantity   int  `json:"quantity"`
	TotalPrice int  `json:"total_price"`
}
