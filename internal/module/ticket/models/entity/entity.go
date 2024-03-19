package entity

import (
	"database/sql"
	"time"
)

type Ticket struct {
	ID        int64        `db:"id"`
	Capacity  int64        `db:"capacity"`
	Region    string       `db:"region"`
	EventDate time.Time    `db:"event_date"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type TicketDetail struct {
	ID        int64        `db:"id"`
	TicketID  int64        `db:"ticket_id"`
	Level     string       `db:"level"`
	Stock     int64        `db:"stock"`
	BasePrice float64      `db:"base_price"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
