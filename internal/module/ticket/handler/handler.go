package handler

import (
	"strconv"
	"ticket-service/internal/module/ticket/models/request"
	"ticket-service/internal/module/ticket/usecases"
	"ticket-service/internal/pkg/errors"
	"ticket-service/internal/pkg/helpers"
	"ticket-service/internal/pkg/log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	Log       log.Logger
	Validator *validator.Validate
	Usecase   usecases.Usecases
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

	// call usecase
	tickets, totalItem, totalPage, err := h.Usecase.ShowTickets(c.Context(), req.Page, req.Size)
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
