package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"ticket-service/internal/module/ticket/models/request"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/module/ticket/repositories"
	"ticket-service/internal/pkg/errors"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gorules/zen-go"
)

type usecases struct {
	repo                 repositories.Repositories
	publish              message.Publisher
	onlineTicketRulesBre zen.Decision
}

// GetTicketByRegionName implements Usecases.
func (u *usecases) GetTicketByRegionName(ctx context.Context, regionName string) (resp []response.Ticket, err error) {
	// get ticket by region name
	ticket, err := u.repo.FindTicketByRegionName(ctx, regionName)
	if err != nil {
		return nil, err
	}

	// find ticket detail
	ticketDetails, err := u.repo.FindTicketDetailByTicketID(ctx, ticket.ID)
	if err != nil {
		return nil, err
	}

	// mapping response

	for _, ticketDetail := range ticketDetails {
		resp = append(resp, response.Ticket{
			ID:        ticketDetail.ID,
			Region:    ticket.Region,
			EventDate: ticket.EventDate,
			Level:     ticketDetail.Level,
			Price:     ticketDetail.BasePrice,
			Stock:     ticketDetail.Stock,
		})
	}

	return resp, nil
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

	// calculate total ticket

	ticketDetails, err := u.repo.FindTicketDetailByTicketID(ctx, ticketDetail.TicketID)
	if err != nil {
		return err
	}

	ticket, err := u.repo.FindTicketByID(ctx, ticketDetail.TicketID)
	if err != nil {
		return err
	}

	var totalTicketPerVenue int64

	for _, ticketDetail := range ticketDetails {
		totalTicketPerVenue += ticketDetail.Stock
	}

	totalTicketPerVenue = totalTicketPerVenue - totalTicket

	spec := request.TicketSoldOut{
		VenueName: ticket.Region,
		IsSoldOut: totalTicketPerVenue < 1,
	}

	payload, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	// check stock
	if totalTicketPerVenue < 1 {
		fmt.Println("publish ticket sold out", spec, totalTicketPerVenue)
		err = u.publish.Publish("update_ticket_sold_out", message.NewMessage(watermill.NewUUID(), payload))
		if err != nil {
			return err
		}
	}

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
	ShowTickets(ctx context.Context, page int, pageSize int, userID int64) (resp []response.Ticket, totalData int, totalPage int, err error)
	// private
	InquiryTicketAmount(ctx context.Context, ticketID int64, totalTicket int) (resp response.InquiryTicketAmount, err error)
	CheckStockTicket(ctx context.Context, ticketDetailID int) (resp response.StockTicket, err error)
	DecrementTicketStock(ctx context.Context, ticketDetailID int64, totalTicket int64) error
	IncrementTicketStock(ctx context.Context, ticketDetailID int64, totalTicket int64) error
	GetTicketByRegionName(ctx context.Context, regionName string) (resp []response.Ticket, err error)
}

func New(repo repositories.Repositories, pub message.Publisher, onlineTicketRulesBre zen.Decision) Usecases {
	return &usecases{
		repo:                 repo,
		publish:              pub,
		onlineTicketRulesBre: onlineTicketRulesBre,
	}
}

func (u *usecases) ShowTickets(ctx context.Context, page int, pageSize int, userID int64) (r []response.Ticket, totalItem int, totalPage int, err error) {
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

	// get ticket region

	var profile response.Profile

	if userID != 0 {
		profile, err = u.repo.GetProfile(ctx, userID)
		if err != nil {
			return nil, 0, 0, err
		}
	}

	if profile.Region == "" {
		profile.Region = "Online"
	}

	venueResult, err := u.repo.GetTicketOnline(ctx, profile.Region)
	if err != nil {
		return nil, 0, 0, err
	}

	// check online ticket seat

	onlineTicket, err := u.repo.FindTicketByRegionName(ctx, "Online")
	if err != nil {
		return nil, 0, 0, err
	}

	// check online ticket rules

	result, err := u.onlineTicketRulesBre.Evaluate(map[string]any{
		"is_ticket_first_sold_out": venueResult.IsFirstSoldOut,
		"is_ticket_sold_out":       venueResult.IsSoldOut,
		"total_seat":               onlineTicket.Capacity,
	})

	if err != nil {
		return nil, 0, 0, err
	}

	var responseBre response.BreOnlineTicket

	byteRes, err := result.Result.MarshalJSON()
	if err != nil {
		return nil, 0, 0, err
	}

	err = json.Unmarshal(byteRes, &responseBre)
	if err != nil {
		return nil, 0, 0, err
	}

	for _, ticket := range tickets {
		for _, td := range ticketDetails {
			if ticket.ID == td.TicketID {
				if ticket.Region == "online" && td.Level == "online" {
					resp = append(resp, response.Ticket{
						ID:        ticket.ID,
						Stock:     responseBre.Seats,
						Region:    ticket.Region,
						Level:     td.Level,
						EventDate: ticket.EventDate,
						Price:     td.BasePrice,
					})
				}
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
