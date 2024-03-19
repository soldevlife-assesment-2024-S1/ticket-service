package usecases

import (
	"context"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/module/ticket/repositories"
)

type usecases struct {
	repo  repositories.Repositories
}

type Usecases interface {
	// public
	ShowTickets(ctx context.Context) ([]response.Ticket, error)
}

func New(repo repositories.Repositories) Usecases {
	return &usecases{
		repo: repo,
	}
}

func (u *usecases) ShowTickets(ctx context.Context) ([]response.Ticket, error) {
	// get data from redis
	tickets, err := u.repo.GetTicketRedis(ctx)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
