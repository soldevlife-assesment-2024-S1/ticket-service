package middleware_test

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"ticket-service/internal/module/ticket/mocks"
// 	"ticket-service/internal/module/ticket/models/response"
// 	log "ticket-service/internal/pkg/log"
// 	log_internal "ticket-service/internal/pkg/log"
// 	"ticket-service/internal/pkg/middleware"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/valyala/fasthttp"
// )

// var (
// 	m        middleware.Middleware
// 	logTest  log.Logger
// 	mockRepo *mocks.Repositories
// 	app      *fiber.App
// )

// func setup() {
// 	logZap := log_internal.SetupLogger()
// 	log_internal.Init(logZap)
// 	logTest := log_internal.GetLogger()
// 	mockRepo = new(mocks.Repositories)
// 	m = middleware.Middleware{
// 		Log:  logTest,
// 		Repo: mockRepo,
// 	}

// 	app = fiber.New()
// }

// func teardown() {
// 	logTest = nil
// 	mockRepo = nil
// 	app = nil
// 	m = middleware.Middleware{}
// }

// func TestMiddlewareValidateToken(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	t.Run("Success Validate Token", func(t *testing.T) {
// 		// mock data
// 		httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
// 		httpReq.Header.Set("Content-Type", "application/json")
// 		httpReq.Header.Set("Authorization", "Bearer token")
// 		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
// 		ctx.Request().Header.Set("Authorization", "Bearer token")
// 		ctx.Request().SetRequestURI("/api/v1/tickets")
// 		ctx.Request().Header.SetMethod(http.MethodGet)
// 		ctx.Request().Header.SetContentType("application/json")

// 		mockResponse := response.UserServiceValidate{
// 			IsValid:   true,
// 			UserID:    1,
// 			EmailUser: "test@test.com",
// 		}

// 		mockRepo.On("ValidateToken", ctx.Context(), "token").Return(mockResponse, nil)

// 		// call function
// 		err := m.ValidateToken(ctx)

// 		// assert
// 		assert.Nil(t, err)
// 	})
// }
