package dto

type TicketRequest struct {
	MovieName      string  `json:"movie_name" binding:"required"`
	TotalTicket    int     `json:"total_ticket" binding:"required,gt=0"`
	PricePerTicket float64 `json:"price_per_ticket" binding:"required,gt=0"`
}

type TicketResponse struct {
	ID              uint    `json:"id"`
	MovieName       string  `json:"movie_name"`
	TotalTicket     int     `json:"total_ticket"`
	AvailableTicket int     `json:"available_ticket"`
	PricePerTicket  float64 `json:"price_per_ticket"`
}
