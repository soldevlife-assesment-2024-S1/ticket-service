package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"ticket-service/config"
	"ticket-service/internal/module/ticket/handler"
	"ticket-service/internal/module/ticket/repositories"
	"ticket-service/internal/module/ticket/usecases"
	"ticket-service/internal/pkg/database"
	"ticket-service/internal/pkg/gorules"
	"ticket-service/internal/pkg/http"
	"ticket-service/internal/pkg/httpclient"
	log_internal "ticket-service/internal/pkg/log"
	"ticket-service/internal/pkg/messagestream"
	"ticket-service/internal/pkg/middleware"
	"ticket-service/internal/pkg/redis"
	router "ticket-service/internal/route"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

func main() {
	cfg := config.InitConfig()

	app, messageRouters := initService(cfg)

	for _, router := range messageRouters {
		ctx := context.Background()
		go func(router *message.Router) {
			err := router.Run(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}(router)
	}

	go func() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		log.Print("Starting runtime instrumentation:")
		err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
		if err != nil {
			log.Fatal(err)
		}

		<-ctx.Done()
	}()

	// start http server
	http.StartHttpServer(app, cfg.HttpServer.Port)
}

func initService(cfg *config.Config) (*fiber.App, []*message.Router) {
	db := database.GetConnection(&cfg.Database)
	redis := redis.SetupClient(&cfg.Redis)
	// logZap := log_internal.SetupLogger()
	// log_internal.Init(logZap)
	// logger := log_internal.GetLogger()
	logger := log_internal.Setup()
	cb := httpclient.InitCircuitBreaker(&cfg.HttpClient, cfg.HttpClient.Type)
	httpClient := httpclient.InitHttpClient(&cfg.HttpClient, cb)

	// init business rules engine
	pathOnlineTicket := "./assets/online-ticket-weight.json"
	breOnlineTicket, err := gorules.Init(pathOnlineTicket)
	if err != nil {
		logger.Ctx(context.Background()).Fatal(fmt.Sprintf("Failed to init BRE: %v", err))
	}

	ctx := context.Background()
	// init message stream
	amqp := messagestream.NewAmpq(&cfg.MessageStream)

	// Init Subscriber
	subscriber, err := amqp.NewSubscriber()
	if err != nil {
		logger.Ctx(ctx).Fatal(fmt.Sprintf("Failed to create subscriber: %v", err))
	}

	// Init Publisher
	publisher, err := amqp.NewPublisher()
	if err != nil {
		logger.Ctx(ctx).Fatal(fmt.Sprintf("Failed to create publisher: %v", err))
	}

	ticketRepo := repositories.New(db, logger, httpClient, redis, &cfg.UserService, &cfg.RecommendationService)
	ticketUsecase := usecases.New(ticketRepo, publisher, breOnlineTicket)
	middleware := middleware.Middleware{
		Repo: ticketRepo,
	}

	validator := validator.New()
	ticketHandler := handler.TicketHandler{
		Log:       logger,
		Validator: validator,
		Usecase:   ticketUsecase,
		Publish:   publisher,
	}

	var messageRouters []*message.Router

	incrementTicketStock, err := messagestream.NewRouter(publisher, "increment_stock_ticket_poisoned", "increment_stock_ticket_handler", "increment_stock_ticket", subscriber, ticketHandler.IncrementTicketStock)
	if err != nil {
		logger.Ctx(ctx).Error(fmt.Sprintf("Failed to create consume_booking_queue router: %v", err))
	}

	decrementTicketStock, err := messagestream.NewRouter(publisher, "decrement_stock_ticket_poisoned", "decrement_stock_ticket_handler", "decrement_stock_ticket", subscriber, ticketHandler.DecrementTicketStock)
	if err != nil {
		logger.Ctx(ctx).Error(fmt.Sprintf("Failed to create consume_booking_queue router: %v", err))
	}

	messageRouters = append(messageRouters, incrementTicketStock, decrementTicketStock)

	serverHttp := http.SetupHttpEngine()
	conn, serviceName, err := http.InitConn(cfg)
	if err != nil {
		logger.Ctx(ctx).Fatal(fmt.Sprintf("Failed to create gRPC connection to collector: %v", err))
	}

	// setup tracer
	http.InitTracer(conn, serviceName)

	// setup metric
	_, err = http.InitMeterProvider(conn, serviceName)
	if err != nil {
		logger.Ctx(ctx).Fatal(fmt.Sprintf("Failed to create meter provider: %v", err))
	}

	r := router.Initialize(serverHttp, &ticketHandler, &middleware)

	return r, messageRouters

}
