package repositories

import (
	"context"
	"fmt"
	"ticket-service/config"
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

type Repositories interface {
	// http
	ValidateToken(ctx context.Context, token string) (bool, error)
	// redis
	GetTicketRedis(ctx context.Context) ([]response.Ticket, error)
	SetTicketRedis(ctx context.Context, tickets []response.Ticket) error
}

func New(db *sqlx.DB, log log.Logger, httpClient *circuit.HTTPClient, redisClient *redis.Client) Repositories {
	return &repositories{
		db:          db,
		log:         log,
		httpClient:  httpClient,
		redisClient: redisClient,
	}
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
