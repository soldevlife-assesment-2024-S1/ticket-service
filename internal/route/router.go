package router

import (
	"ticket-service/internal/module/ticket/handler"
	"ticket-service/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func Initialize(app *fiber.App, handlerTicket *handler.TicketHandler, m *middleware.Middleware) *fiber.App {

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})

	Api := app.Group("/api")

	// public routes
	v1 := Api.Group("/v1")
	// v1.Get("/tickets", m.ValidateToken, handlerTicket.ShowTickets)
	v1.Get("/tickets", handlerTicket.ShowTickets)

	private := Api.Group("/private")
	private.Get("/ticket/inquiry", handlerTicket.InquiryTicketAmount)
	private.Get("/ticket/stock", handlerTicket.CheckStockTicket)
	private.Get("ticket", handlerTicket.GetTicketByRegionName)

	return app

}
