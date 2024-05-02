package handler

import (
	"context"
	"strconv"
	"ticket-service/internal/module/ticket/models/request"
	"ticket-service/internal/module/ticket/usecases"
	"ticket-service/internal/pkg/errors"
	"ticket-service/internal/pkg/helpers"
	"ticket-service/internal/pkg/log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	Log       log.Logger
	Validator *validator.Validate
	Usecase   usecases.Usecases
	Publish   message.Publisher
}

func (h *TicketHandler) ShowTickets(c *fiber.Ctx) error {
	var req request.Pagination
	if err := c.QueryParser(&req); err != nil {
		return helpers.RespError(c, h.Log, errors.BadRequest("Bad Request, Invalid Query Params"))
	}

	// validate request
	if err := h.Validator.Struct(req); err != nil {
		return helpers.RespError(c, h.Log, errors.BadRequest(err.Error()))
	}

	userID := c.Locals("user_id").(int64)

	// call usecase
	tickets, totalItem, totalPage, err := h.Usecase.ShowTickets(c.Context(), req.Page, req.Size, userID)
	if err != nil {
		return helpers.RespError(c, h.Log, err)
	}

	meta := helpers.MetaPaginationResponse{
		Code:      200,
		Message:   "Show Tickets Success",
		Page:      req.Page,
		Size:      req.Size,
		TotalPage: totalPage,
		TotalData: totalItem,
	}

	// response
	return helpers.RespPagination(c, h.Log, tickets, meta, "Show Tickets Success")
}

// private

func (h *TicketHandler) InquiryTicketAmount(c *fiber.Ctx) error {
	ticketDetailID := c.Query("ticket_detail_id")
	totalTicket := c.Query("total_ticket")

	intTicketDetailID, err := strconv.Atoi(ticketDetailID)
	if err != nil {
		return helpers.RespError(c, h.Log, errors.BadRequest("Invalid Ticket Detail ID"))
	}

	int64TicketDetailID := int64(intTicketDetailID)

	intTotalTicket, err := strconv.Atoi(totalTicket)
	if err != nil {
		return helpers.RespError(c, h.Log, errors.BadRequest("Invalid Total Ticket"))
	}
	// call usecase
	amount, err := h.Usecase.InquiryTicketAmount(c.Context(), int64TicketDetailID, intTotalTicket)
	if err != nil {
		return helpers.RespError(c, h.Log, err)
	}

	// response
	return helpers.RespSuccess(c, h.Log, amount, "Inquiry Ticket Amount Success")
}

func (h *TicketHandler) CheckStockTicket(c *fiber.Ctx) error {

	ticketDetailID := c.Query("ticket_detail_id")

	intTicketDetailID, err := strconv.Atoi(ticketDetailID)
	if err != nil {
		return helpers.RespError(c, h.Log, errors.BadRequest("Invalid Ticket Detail ID"))
	}

	// call usecase
	amount, err := h.Usecase.CheckStockTicket(c.Context(), intTicketDetailID)
	if err != nil {
		return helpers.RespError(c, h.Log, err)
	}

	// response
	return helpers.RespSuccess(c, h.Log, amount, "Check Stock Ticket Success")
}

func (h *TicketHandler) DecrementTicketStock(msg *message.Message) error {
	msg.Ack()
	req := request.DecrementTicketStock{}

	if err := json.Unmarshal(msg.Payload, &req); err != nil {
		h.Log.Error(msg.Context(), "Failed to unmarshal data", err)

		// publish to poison queue
		reqPoisoned := request.PoisonedQueue{
			TopicTarget: "decrement_stock_ticket",
			ErrorMsg:    err.Error(),
			Payload:     msg.Payload,
		}

		jsonPayload, _ := json.Marshal(reqPoisoned)

		err = h.Publish.Publish("poisoned_queue", message.NewMessage(watermill.NewUUID(), jsonPayload))
		if err != nil {
			h.Log.Error(msg.Context(), "Failed to publish to poison queue", err)
		}

		return err
	}

	ctx := context.Background()

	// call usecase
	err := h.Usecase.DecrementTicketStock(ctx, req.TicketDetailID, req.TotalTickets)

	if err != nil {
		h.Log.Error(msg.Context(), "Failed to decrement ticket stock", err)

		// publish to poison queue
		reqPoisoned := request.PoisonedQueue{
			TopicTarget: "decrement_stock_ticket",
			ErrorMsg:    err.Error(),
			Payload:     msg.Payload,
		}

		jsonPayload, _ := json.Marshal(reqPoisoned)

		err = h.Publish.Publish("poisoned_queue", message.NewMessage(watermill.NewUUID(), jsonPayload))
		if err != nil {
			h.Log.Error(msg.Context(), "Failed to publish to poison queue", err)
		}

		return err
	}

	return nil
}

func (h *TicketHandler) IncrementTicketStock(msg *message.Message) error {
	msg.Ack()
	req := request.IncrementTicketStock{}

	if err := json.Unmarshal(msg.Payload, &req); err != nil {
		h.Log.Error(msg.Context(), "Failed to unmarshal data", err)

		// publish to poison queue
		reqPoisoned := request.PoisonedQueue{
			TopicTarget: "increment_stock_ticket",
			ErrorMsg:    err.Error(),
			Payload:     msg.Payload,
		}

		jsonPayload, _ := json.Marshal(reqPoisoned)

		err = h.Publish.Publish("poisoned_queue", message.NewMessage(watermill.NewUUID(), jsonPayload))
		if err != nil {
			h.Log.Error(msg.Context(), "Failed to publish to poison queue", err)
		}

		return err
	}
	ctx := context.Background()

	// call usecase
	err := h.Usecase.IncrementTicketStock(ctx, req.TicketDetailID, req.TotalTickets)

	if err != nil {
		h.Log.Error(msg.Context(), "Failed to increment ticket stock", err)

		// publish to poison queue
		reqPoisoned := request.PoisonedQueue{
			TopicTarget: "increment_stock_ticket",
			ErrorMsg:    err.Error(),
			Payload:     msg.Payload,
		}

		jsonPayload, _ := json.Marshal(reqPoisoned)

		err = h.Publish.Publish("poisoned_queue", message.NewMessage(watermill.NewUUID(), jsonPayload))
		if err != nil {
			h.Log.Error(msg.Context(), "Failed to publish to poison queue", err)
		}

		return err
	}

	return nil
}

func (h *TicketHandler) GetTicketByRegionName(c *fiber.Ctx) error {
	regionName := c.Query("region_name")

	// call usecase
	tickets, err := h.Usecase.GetTicketByRegionName(c.Context(), regionName)
	if err != nil {
		return helpers.RespError(c, h.Log, err)
	}

	// response
	return helpers.RespSuccess(c, h.Log, tickets, "Get Ticket By Region Name Success")

}
