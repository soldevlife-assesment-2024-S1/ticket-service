package request

type Pagination struct {
	Page int `json:"page" form:"page" required:"true" validate:"required,numeric"`
	Size int `json:"size" form:"size" required:"true" validate:"required,numeric"`
}

type InquiryTicketAmount struct {
	TicketID    int64 `json:"ticket_id" form:"ticket_id" required:"true" validate:"required"`
	TotalTicket int   `json:"total_ticket" form:"total_ticket" required:"true" validate:"required,numeric"`
}

type CheckStockTicket struct {
	TicketDetailID string `form:"ticket_detail_id"`
}

type DecrementTicketStock struct {
	TicketDetailID int64 `json:"ticket_detail_id" form:"ticket_detail_id" validate:"required"`
	TotalTickets   int64 `json:"total_tickets" form:"total_tickets" validate:"required"`
}

type IncrementTicketStock struct {
	TicketDetailID int64 `json:"ticket_detail_id" form:"ticket_detail_id" validate:"required"`
	TotalTickets   int64 `json:"total_tickets" form:"total_tickets" validate:"required"`
}

type PoisonedQueue struct {
	TopicTarget string      `json:"topic_target" validate:"required"`
	ErrorMsg    string      `json:"error_msg" validate:"required"`
	Payload     interface{} `json:"payload" validate:"required"`
}
