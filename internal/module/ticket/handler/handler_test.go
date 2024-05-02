package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"ticket-service/internal/module/ticket/handler"
	"ticket-service/internal/module/ticket/mocks"
	"ticket-service/internal/module/ticket/models/request"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/pkg/log"
	log_internal "ticket-service/internal/pkg/log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

var (
	h             *handler.TicketHandler
	uc            *mocks.Usecases
	p             message.Publisher
	logMock       log.Logger
	app           *fiber.App
	validatorTest *validator.Validate
)

type mockPublisher struct{}

// Close implements message.Publisher.
func (m *mockPublisher) Close() error {
	return nil
}

// Publish implements message.Publisher.
func (m *mockPublisher) Publish(topic string, messages ...*message.Message) error {
	return nil
}

func NewMockPublisher() message.Publisher {
	return &mockPublisher{}
}

func setup() {
	uc = new(mocks.Usecases)
	p = NewMockPublisher()
	logZap := log_internal.SetupLogger()
	log_internal.Init(logZap)
	logMock := log_internal.GetLogger()
	validatorTest = validator.New()
	h = &handler.TicketHandler{
		Usecase:   uc,
		Publish:   p,
		Log:       logMock,
		Validator: validatorTest,
	}

	app = fiber.New()
}

func teardown() {
	h = nil
	uc = nil
	p = nil
	logMock = nil
}

func TestShowTickets(t *testing.T) {
	setup()
	defer teardown()

	t.Run("success", func(t *testing.T) {
		// mock data
		payload := request.Pagination{
			Page: 1,
			Size: 10,
		}

		httpReq := httptest.NewRequest(http.MethodGet, "/tickets", nil)
		httpReq.Header.Set("Content-Type", "application/json")

		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		ctx.Request().SetRequestURI("api/v1/tickets")
		ctx.Request().Header.SetMethod(http.MethodGet)
		ctx.Request().Header.SetContentType("application/json")
		ctx.Request().URI().QueryArgs().Add("page", strconv.Itoa(payload.Page))
		ctx.Request().URI().QueryArgs().Add("size", strconv.Itoa(payload.Size))
		ctx.Locals("user_id", int64(1))

		// mock usecase
		uc.On("ShowTickets", ctx.Context(), payload.Page, payload.Size, int64(1)).Return(nil, 0, 0, nil)

		// call handler
		err := h.ShowTickets(ctx)
		fmt.Println("error ticket", err)

		// assert
		assert.Nil(t, err)

	})
}

func TestInquiryTicketAmount(t *testing.T) {
	setup()
	defer teardown()

	t.Run("success", func(t *testing.T) {
		// mock data
		httpReq := httptest.NewRequest(http.MethodGet, "/ticket/inquiry", nil)
		httpReq.Header.Set("Content-Type", "application/json")

		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		ctx.Request().SetRequestURI("api/private/ticket/inquiry")
		ctx.Request().Header.SetMethod(http.MethodGet)
		ctx.Request().Header.SetContentType("application/json")

		// mock usecase
		uc.On("InquiryTicketAmount", ctx.Context()).Return(nil)

		// call handler
		err := h.InquiryTicketAmount(ctx)
		fmt.Println("error ticket", err)

		// assert
		assert.Nil(t, err)

	})
}

func TestCheckStockTicket(t *testing.T) {
	setup()
	defer teardown()

	t.Run("success", func(t *testing.T) {
		// mock data
		httpReq := httptest.NewRequest(http.MethodGet, "/ticket/stock", nil)
		httpReq.Header.Set("Content-Type", "application/json")

		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		ctx.Request().SetRequestURI("api/private/ticket/stock")
		ctx.Request().Header.SetMethod(http.MethodGet)
		ctx.Request().Header.SetContentType("application/json")

		// mock usecase
		uc.On("CheckStockTicket", ctx.Context()).Return(nil)

		// call handler
		err := h.CheckStockTicket(ctx)
		fmt.Println("error ticket", err)

		// assert
		assert.Nil(t, err)

	})
}

func TestGetTicketByRegionName(t *testing.T) {
	setup()
	defer teardown()

	t.Run("success", func(t *testing.T) {
		// mock data

		responses := []response.Ticket{
			{
				ID:     1,
				Level:  "Silver",
				Price:  1000,
				Stock:  10,
				Region: "Asean",
			},
			{
				ID:     2,
				Level:  "Gold",
				Price:  2000,
				Stock:  20,
				Region: "Asean",
			},
		}
		httpReq := httptest.NewRequest(http.MethodGet, "/api/private/ticket", nil)
		httpReq.Header.Set("Content-Type", "application/json")

		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		ctx.Request().SetRequestURI("/api/private/ticket")
		ctx.Request().Header.SetMethod(http.MethodGet)
		ctx.Request().Header.SetContentType("application/json")
		ctx.Request().URI().QueryArgs().Add("region_name", "Asean")

		// mock usecase
		uc.On("GetTicketByRegionName", ctx.Context(), "Asean").Return(responses, nil)

		// call handler
		err := h.GetTicketByRegionName(ctx)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, ctx.Response().StatusCode())

	})
}
