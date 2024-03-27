package response

import "time"

type UserServiceValidate struct {
	IsValid bool `json:"is_valid"`
}

type Ticket struct {
	ID        int64     `json:"id"`
	Stock     int64     `json:"stock"`
	Region    string    `json:"region"`
	Level     string    `json:"level"`
	EventDate time.Time `json:"event_date"`
	Price     float64   `json:"price"`
}

type InquiryTicketAmount struct {
	TotalTicket int     `json:"total_ticket"`
	TotalAmount float64 `json:"total_amount"`
}

type StockTicket struct {
	Stock int64 `json:"stock"`
}
