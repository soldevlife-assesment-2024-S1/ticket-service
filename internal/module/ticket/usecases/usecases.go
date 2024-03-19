package usecases

import (
	"context"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/module/ticket/repositories"
)

type usecases struct {
	repo repositories.Repositories
}

type Usecases interface {
	// public
	ShowTickets(ctx context.Context, page int, pageSize int) (resp []response.Ticket, totalData int, totalPage int, err error)
}

func New(repo repositories.Repositories) Usecases {
	return &usecases{
		repo: repo,
	}
}

func (u *usecases) ShowTickets(ctx context.Context, page int, pageSize int) (r []response.Ticket, totalItem int, totalPage int, err error) {
	// // get data from redis
	// tickets, err := u.repo.GetTicketRedis(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	var resp []response.Ticket

	// get data from database
	tickets, totalItem, totalPage, err := u.repo.FindTickets(ctx, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}

	ticketDetails, totalItem, totalPage, err := u.repo.FindTicketDetails(ctx, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}

	for _, ticket := range tickets {
		for _, td := range ticketDetails {
			if ticket.ID == td.TicketID {
				resp = append(resp, response.Ticket{
					ID:        ticket.ID,
					Capacity:  ticket.Capacity,
					Region:    ticket.Region,
					Level:     ticket.Level,
					EventDate: ticket.EventDate,
					Price:     td.BasePrice,
				})
			}
		}
	}

	return resp, totalItem, totalPage, nil
}
