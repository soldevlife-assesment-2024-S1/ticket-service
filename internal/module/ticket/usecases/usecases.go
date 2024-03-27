package usecases

import (
	"context"
	"fmt"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/module/ticket/repositories"
	"ticket-service/internal/pkg/errors"
)

type usecases struct {
	repo repositories.Repositories
}

// DecrementTicketStock implements Usecases.
func (u *usecases) DecrementTicketStock(ctx context.Context, ticketDetailID int64, totalTicket int64) error {

	// get ticket detail
	ticketDetail, err := u.repo.FindTicketDetail(ctx, ticketDetailID)
	if err != nil {
		return err
	}

	// check stock
	if ticketDetail.Stock < totalTicket {
		return errors.BadRequest("stock not enough")
	}

	// decrement stock
	ticketDetail.Stock -= totalTicket

	fmt.Println("stock", ticketDetail.Stock)

	// update stock
	err = u.repo.UpdateTicketDetail(ctx, ticketDetail)
	if err != nil {
		return err
	}

	return nil
}

// IncrementTicketStock implements Usecases.
func (u *usecases) IncrementTicketStock(ctx context.Context, ticketDetailID int64, totalTicket int64) error {
	// get ticket detail
	ticketDetail, err := u.repo.FindTicketDetail(ctx, ticketDetailID)
	if err != nil {
		return err
	}

	// increment stock
	ticketDetail.Stock += totalTicket

	// update stock
	err = u.repo.UpdateTicketDetail(ctx, ticketDetail)
	if err != nil {
		return err
	}

	return nil
}

// CheckStockTicket implements Usecases.
func (u *usecases) CheckStockTicket(ctx context.Context, ticketDetailID int) (resp response.StockTicket, err error) {
	// get ticket detail
	ticketDetailID64 := int64(ticketDetailID)
	ticketDetail, err := u.repo.FindTicketDetail(ctx, ticketDetailID64)
	if err != nil {
		return response.StockTicket{}, err
	}

	// check stock
	if ticketDetail.Stock == 0 {
		return response.StockTicket{
			Stock: 0,
		}, nil
	}

	resp = response.StockTicket{
		Stock: ticketDetail.Stock,
	}

	return resp, nil
}

// InquiryTicketAmount implements Usecases.
func (u *usecases) InquiryTicketAmount(ctx context.Context, ticketID int64, totalTicket int) (resp response.InquiryTicketAmount, err error) {
	// get ticket details
	ticketDetails, err := u.repo.FindTicketDetail(ctx, ticketID)
	if err != nil {
		return response.InquiryTicketAmount{}, err
	}

	// calculate total amount
	totalAmount := ticketDetails.BasePrice * float64(totalTicket)

	return response.InquiryTicketAmount{
		TotalTicket: totalTicket,
		TotalAmount: totalAmount,
	}, nil
}

type Usecases interface {
	// public
	ShowTickets(ctx context.Context, page int, pageSize int) (resp []response.Ticket, totalData int, totalPage int, err error)
	// private
	InquiryTicketAmount(ctx context.Context, ticketID int64, totalTicket int) (resp response.InquiryTicketAmount, err error)
	CheckStockTicket(ctx context.Context, ticketDetailID int) (resp response.StockTicket, err error)
	DecrementTicketStock(ctx context.Context, ticketDetailID int64, totalTicket int64) error
	IncrementTicketStock(ctx context.Context, ticketDetailID int64, totalTicket int64) error
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
	tickets, _, _, err := u.repo.FindTickets(ctx, page, pageSize)
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
					Stock:     td.Stock,
					Region:    ticket.Region,
					Level:     td.Level,
					EventDate: ticket.EventDate,
					Price:     td.BasePrice,
				})
			}
		}
	}

	return resp, totalItem, totalPage, nil
}
