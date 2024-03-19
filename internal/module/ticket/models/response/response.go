package response

import "time"

type UserServiceValidate struct {
	IsValid bool `json:"is_valid"`
}

type Ticket struct {
	ID        int       `json:"id"`
	Capacity  int       `json:"capacity"`
	Region    int       `json:"region"`
	Level     int       `json:"level"`
	EventDate time.Time `json:"event_date"`
}
