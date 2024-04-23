package response

import "time"

type BaseResponse struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}
type UserServiceValidate struct {
	IsValid   bool   `json:"is_valid"`
	UserID    int64  `json:"user_id"`
	EmailUser string `json:"email_user"`
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

type BreOnlineTicket struct {
	Seats int64 `json:"seats"`
}

type OnlineTicket struct {
	IsSoldOut      bool `json:"is_sold_out"`
	IsFirstSoldOut bool `json:"is_first_sold_out"`
}

type Profile struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Address        string `json:"address"`
	District       string `json:"district"`
	City           string `json:"city"`
	State          string `json:"state"`
	Country        string `json:"country"`
	Region         string `json:"region"`
	Phone          string `json:"phone"`
	PersonalID     string `json:"personal_id"`
	TypePersonalID string `json:"type_personal_id"`
}
