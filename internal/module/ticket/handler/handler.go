package handler

import (
	"ticket-service/internal/module/ticket/usecases"
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
	// call usecase
	tickets, err := h.Usecase.ShowTickets(c.Context())
	if err != nil {
		return helpers.RespError(c, h.Log, err)
	}

	// response
	return helpers.RespSuccess(c, h.Log, tickets, "Success")
}
