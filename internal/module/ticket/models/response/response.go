package response

import "time"

type UserServiceValidate struct {
	IsValid bool `json:"is_valid"`
}

type Ticket struct {
	ID        int64     `json:"id"`
	Capacity  int64     `json:"capacity"`
	Region    string    `json:"region"`
	Level     string    `json:"level"`
	EventDate time.Time `json:"event_date"`
	Price     float64   `json:"price"`
}
