package handler

import (
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
