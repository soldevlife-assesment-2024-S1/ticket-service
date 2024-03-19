package entity

import "time"

type Ticket struct {
	ID        int `db:"id"`
	Capacity  int `db:"capacity"`
	Region    int `db:"region"`
	Level     int `db:"level"`
	EventDate time.Time `db:"event_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}


type TicketDetails struct {
	ID        int       `db:"id"`
	TicketID  int       `db:"ticket_id"`
	BasePrice int       `db:"base_price"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}