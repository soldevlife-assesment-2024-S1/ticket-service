package repositories

import (
	"context"
	"fmt"
	"ticket-service/config"
	"ticket-service/internal/module/ticket/models/entity"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/pkg/errors"
	"ticket-service/internal/pkg/log"

	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	circuit "github.com/rubyist/circuitbreaker"
)

type repositories struct {
	db             *sqlx.DB
	log            log.Logger
	httpClient     *circuit.HTTPClient
	cfgUserService *config.UserService
	redisClient    *redis.Client
}

// FindTicketDetail implements Repositories.
func (r *repositories) FindTicketDetail(ctx context.Context, ticketID int64) (entity.TicketDetail, error) {
	query := fmt.Sprintf("SELECT * FROM ticket_details WHERE id = %d", ticketID)
	var ticketDetail entity.TicketDetail
	if err := r.db.GetContext(ctx, &ticketDetail, query); err != nil {
		r.log.Error(ctx, "From Repositories: Failed to execute query", err)
		return entity.TicketDetail{}, err
	}

	return ticketDetail, nil
}

type Repositories interface {
	// http
	ValidateToken(ctx context.Context, token string) (bool, error)
	// redis
	GetTicketRedis(ctx context.Context) ([]response.Ticket, error)
	SetTicketRedis(ctx context.Context, tickets []response.Ticket) error
	FindTickets(ctx context.Context, page int, pageSize int) (tickets []entity.Ticket, totalCount int, totalPage int, err error)
	FindTicketDetails(ctx context.Context, page int, pageSize int) (ticketDetails []entity.TicketDetail, totalCount int, totalPage int, err error)
	FindTicketDetail(ctx context.Context, ticketID int64) (entity.TicketDetail, error)
}

func New(db *sqlx.DB, log log.Logger, httpClient *circuit.HTTPClient, redisClient *redis.Client) Repositories {
	return &repositories{
		db:          db,
		log:         log,
		httpClient:  httpClient,
		redisClient: redisClient,
	}
}

func (r *repositories) FindTickets(ctx context.Context, page int, pageSize int) (tickets []entity.Ticket, totalCount int, totalPage int, err error) {
	// calculate offset and limit
	offset := (page - 1) * pageSize
	limit := pageSize

	// query with pagination
	query := fmt.Sprintf("SELECT * FROM tickets LIMIT %d OFFSET %d", limit, offset)

	// execute query
	ticketsErrCh := make(chan error)
	go func() {
		if err := r.db.SelectContext(ctx, &tickets, query); err != nil {
			r.log.Error(ctx, "From Repositories: Failed to execute query", err)
			ticketsErrCh <- err
		}
		close(ticketsErrCh)
	}()

	// get total item count
	totalCountCh := make(chan int)
	go func() {
		var totalCount int
		if err := r.db.GetContext(ctx, &totalCount, "SELECT COUNT(*) FROM tickets"); err != nil {
			r.log.Error(ctx, "From Repositories: Failed to get total item count", err)
			totalCountCh <- 0
		}
		totalCountCh <- totalCount
		close(totalCountCh)
	}()

	// calculate total page
	totalCount = <-totalCountCh
	totalPage = totalCount / pageSize
	if totalCount%pageSize != 0 {
		totalPage++
	}

	// check for errors
	if err := <-ticketsErrCh; err != nil {
		return nil, 0, 0, err
	}

	return tickets, totalCount, totalPage, nil
}

func (r *repositories) FindTicketDetails(ctx context.Context, page int, pageSize int) (ticketDetails []entity.TicketDetail, totalCount int, totalPage int, err error) {
	// calculate offset and limit
	offset := (page - 1) * pageSize
	limit := pageSize

	// query with pagination
	query := fmt.Sprintf("SELECT * FROM ticket_details LIMIT %d OFFSET %d", limit, offset)

	// execute query
	ticketDetailsErrCh := make(chan error)
	go func() {
		if err := r.db.SelectContext(ctx, &ticketDetails, query); err != nil {
			r.log.Error(ctx, "From Repositories: Failed to execute query", err)
			ticketDetailsErrCh <- err
		}
		close(ticketDetailsErrCh)
	}()

	// get total item count
	totalCountCh := make(chan int)
	go func() {
		var totalCount int
		if err := r.db.GetContext(ctx, &totalCount, "SELECT COUNT(*) FROM ticket_details"); err != nil {
			r.log.Error(ctx, "From Repositories: Failed to get total item count", err)
			totalCountCh <- 0
		}
		totalCountCh <- totalCount
		close(totalCountCh)
	}()

	// calculate total page
	totalCount = <-totalCountCh
	totalPage = totalCount / pageSize
	if totalCount%pageSize != 0 {
		totalPage++
	}

	// check for errors
	if err := <-ticketDetailsErrCh; err != nil {
		return nil, 0, 0, err
	}

	return ticketDetails, totalCount, totalPage, nil
}

func (r *repositories) SetTicketRedis(ctx context.Context, tickets []response.Ticket) error {
	// set data to redis
	val, err := json.Marshal(tickets)
	if err != nil {
		r.log.Error(ctx, "From Repositories: Failed to marshal data", err)
		return errors.BadRequest("Failed to marshal data")
	}

	if err := r.redisClient.Set(ctx, "tickets", val, 0).Err(); err != nil {
		r.log.Error(ctx, "From Repositories: Failed to set data to redis", err)
		return err
	}

	return nil
}

func (r *repositories) GetTicketRedis(ctx context.Context) ([]response.Ticket, error) {
	// get data from redis
	val, err := r.redisClient.Get(ctx, "tickets").Result()
	if err != nil {
		r.log.Error(ctx, "From Repositories: Failed to get data from redis", err)
		return nil, err
	}

	// parse response
	var tickets []response.Ticket
	if err := json.Unmarshal([]byte(val), &tickets); err != nil {
		r.log.Error(ctx, "From Repositories: Failed to unmarshal data", err)
		return nil, err
	}

	return tickets, nil
}

func (r *repositories) ValidateToken(ctx context.Context, token string) (bool, error) {
	// http call to user service
	url := fmt.Sprintf("http://%s:%s/api/private/token/validate?token=%s", r.cfgUserService.Host, r.cfgUserService.Port, token)
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		r.log.Error(ctx, "Invalid token", resp.StatusCode)
		return false, errors.BadRequest("Invalid token")
	}

	// parse response
	var respData response.UserServiceValidate

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respData); err != nil {
		return false, err
	}

	if !respData.IsValid {
		r.log.Error(ctx, "Invalid token", resp.StatusCode)
		return false, errors.BadRequest("Invalid token")
	}

	// validate token
	return true, nil
}
