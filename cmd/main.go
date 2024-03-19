package main

import (
	"ticket-service/config"
	"ticket-service/internal/module/ticket/handler"
	"ticket-service/internal/module/ticket/repositories"
	"ticket-service/internal/module/ticket/usecases"
	"ticket-service/internal/pkg/database"
	"ticket-service/internal/pkg/http"
	"ticket-service/internal/pkg/httpclient"
	"ticket-service/internal/pkg/log"
	"ticket-service/internal/pkg/middleware"
	"ticket-service/internal/pkg/redis"
	router "ticket-service/internal/route"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.InitConfig()

	app := initService(cfg)

	// start http server
	http.StartHttpServer(app, cfg.HttpServer.Port)
}

func initService(cfg *config.Config) *fiber.App {
	db := database.GetConnection(&cfg.Database)
	redis := redis.SetupClient(&cfg.Redis)
	logZap := log.SetupLogger()
	log.Init(logZap)
	logger := log.GetLogger()
	cb := httpclient.InitCircuitBreaker(&cfg.HttpClient, cfg.HttpClient.Type)
	httpClient := httpclient.InitHttpClient(&cfg.HttpClient, cb)

	ticketRepo := repositories.New(db, logger, httpClient, redis)
	ticketUsecase := usecases.New(ticketRepo)
	middleware := middleware.Middleware{
		Repo: ticketRepo,
	}

	validator := validator.New()
	userHandler := handler.TicketHandler{
		Log:       logger,
		Validator: validator,
		Usecase:   ticketUsecase,
	}

	serverHttp := http.SetupHttpEngine()

	r := router.Initialize(serverHttp, &userHandler, &middleware)

	return r

}
