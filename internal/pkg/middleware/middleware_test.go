package middleware_test

import (
	"ticket-service/internal/module/ticket/mocks"
	log "ticket-service/internal/pkg/log"
	"ticket-service/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

var (
	m        middleware.Middleware
	logTest  log.Logger
	mockRepo *mocks.Repositories
	app      *fiber.App
)

func setup() {
	logTest = log.SetupLogger()
	mockRepo = new(mocks.Repositories)
	app = fiber.New()
	m = middleware.Middleware{
		Log:  logTest,
		Repo: mockRepo,
	}
}

func teardown() {
	logTest = nil
	mockRepo = nil
	app = nil
	m = middleware.Middleware{}
}

// func TestMiddleware_ValidateToken(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	t.Run("Success Validate Token", func(t *testing.T) {
// 		// mock data
// 		httpReq := httptest.NewRequest(http.MethodGet, "/tickets", nil)
// 		httpReq.Header.Set("Content-Type", "application/json")
// 		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
// 		ctx.Request().Header.Set("Authorization", "Bearer token")
// 		ctx.Request().SetRequestURI("api/v1/tickets")
// 		ctx.Request().Header.SetMethod(http.MethodGet)
// 		ctx.Request().Header.SetContentType("application/json")

// 		mockResponse := response.UserServiceValidate{
// 			IsValid:   true,
// 			UserID:    1,
// 			EmailUser: "test@test.com",
// 		}

// 		mockRepo.On("ValidateToken", ctx.Context(), "token").Return(mockResponse, nil)

// 		// call function
// 		m.ValidateToken(ctx)
// 	})
// }
