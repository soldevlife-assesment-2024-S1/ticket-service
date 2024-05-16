package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"ticket-service/config"
	"ticket-service/internal/module/ticket/models/entity"
	"ticket-service/internal/module/ticket/models/response"
	"ticket-service/internal/pkg/errors"

	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	circuit "github.com/rubyist/circuitbreaker"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type repositories struct {
	db                       *sqlx.DB
	log                      *otelzap.Logger
	httpClient               *circuit.HTTPClient
	cfgUserService           *config.UserService
	cfgRecommendationService *config.RecommendationServiceConfig
	redisClient              *redis.Client
}

// FindTicketByID implements Repositories.
func (r *repositories) FindTicketByID(ctx context.Context, ticketID int64) (entity.Ticket, error) {
	query := fmt.Sprintf("SELECT * FROM tickets WHERE id = %d", ticketID)
	var ticket entity.Ticket
	if err := r.db.GetContext(ctx, &ticket, query); err != nil {
		r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket by id, %v", err))
		return entity.Ticket{}, err
	}

	return ticket, nil
}

// FindTicketDetailByTicketID implements Repositories.
func (r *repositories) FindTicketDetailByTicketID(ctx context.Context, ticketID int64) ([]entity.TicketDetail, error) {
	query := fmt.Sprintf("SELECT * FROM ticket_details WHERE ticket_id = %d", ticketID)
	var ticketDetails []entity.TicketDetail
	if err := r.db.SelectContext(ctx, &ticketDetails, query); err != nil {
		r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket by id, %v", err))
		return ticketDetails, err
	}

	return ticketDetails, nil
}

// FindTicketByRegionName implements Repositories.
func (r *repositories) FindTicketByRegionName(ctx context.Context, regionName string) (entity.Ticket, error) {
	regionName = strings.ToLower(regionName)
	query := fmt.Sprintf("SELECT * FROM tickets WHERE region = '%s'", regionName)
	var ticket entity.Ticket
	if err := r.db.GetContext(ctx, &ticket, query); err != nil {
		r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket by id, %v", err))
		return ticket, err
	}

	return ticket, nil
}

// UpsertTicketDetail implements Repositories.
func (r *repositories) UpdateTicketDetail(ctx context.Context, ticketDetail entity.TicketDetail) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.InternalServerError("error starting transaction")
	}

	// Lock the rows for update
	query := `SELECT * FROM ticket_details WHERE id = $1 FOR UPDATE`
	var existingTicketDetail entity.TicketDetail
	err = r.db.GetContext(ctx, &existingTicketDetail, query, ticketDetail.ID)
	if err != nil && err != sql.ErrNoRows {
		err = tx.Rollback()
		if err != nil {
			return errors.InternalServerError("error rolling back transaction")
		}
		return errors.InternalServerError("error locking rows")
	}

	// Update existing ticket detail
	queryUpdate := `UPDATE ticket_details SET stock = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.ExecContext(ctx, queryUpdate, ticketDetail.Stock, ticketDetail.ID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return errors.InternalServerError("error rolling back transaction")
		}
		return errors.InternalServerError("error upserting ticket detail")
	}

	err = tx.Commit()
	if err != nil {
		return errors.InternalServerError("error committing transaction")
	}

	return nil
}

// FindTicketDetail implements Repositories.
func (r *repositories) FindTicketDetail(ctx context.Context, ticketID int64) (entity.TicketDetail, error) {
	query := fmt.Sprintf("SELECT * FROM ticket_details WHERE id = %d", ticketID)
	var ticketDetail entity.TicketDetail
	if err := r.db.GetContext(ctx, &ticketDetail, query); err != nil {
		r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket by id, %v", err))
		return entity.TicketDetail{}, err
	}

	return ticketDetail, nil
}

type Repositories interface {
	// http
	ValidateToken(ctx context.Context, token string) (response.UserServiceValidate, error)
	GetTicketOnline(ctx context.Context, regionName string) (response.OnlineTicket, error)
	GetProfile(ctx context.Context, userID int64) (response.Profile, error)
	// redis
	// GetTicketRedis(ctx context.Context) ([]response.Ticket, error)
	// SetTicketRedis(ctx context.Context, tickets []response.Ticket) error
	// db
	FindTickets(ctx context.Context, page int, pageSize int) (tickets []entity.Ticket, totalCount int, totalPage int, err error)
	FindTicketByID(ctx context.Context, ticketID int64) (entity.Ticket, error)
	FindTicketDetails(ctx context.Context, page int, pageSize int) (ticketDetails []entity.TicketDetail, totalCount int, totalPage int, err error)
	FindTicketDetail(ctx context.Context, ticketID int64) (entity.TicketDetail, error)
	UpdateTicketDetail(ctx context.Context, ticketDetail entity.TicketDetail) error
	FindTicketByRegionName(ctx context.Context, regionName string) (entity.Ticket, error)
	FindTicketDetailByTicketID(ctx context.Context, ticketID int64) ([]entity.TicketDetail, error)
}

func New(db *sqlx.DB, log *otelzap.Logger, httpClient *circuit.HTTPClient, redisClient *redis.Client, cfgUserService *config.UserService, cfgRecommendationService *config.RecommendationServiceConfig) Repositories {
	return &repositories{
		db:                       db,
		log:                      log,
		httpClient:               httpClient,
		redisClient:              redisClient,
		cfgUserService:           cfgUserService,
		cfgRecommendationService: cfgRecommendationService,
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
			r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket by id, %v", err))
			ticketsErrCh <- err
		}
		close(ticketsErrCh)
	}()

	// get total item count
	totalCountCh := make(chan int)
	go func() {
		var totalCount int
		if err := r.db.GetContext(ctx, &totalCount, "SELECT COUNT(*) FROM tickets"); err != nil {
			r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket by id, %v", err))
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
			r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket by id, %v", err))
			ticketDetailsErrCh <- err
		}
		close(ticketDetailsErrCh)
	}()

	// get total item count
	totalCountCh := make(chan int)
	go func() {
		var totalCount int
		if err := r.db.GetContext(ctx, &totalCount, "SELECT COUNT(*) FROM ticket_details"); err != nil {
			r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket by id, %v", err))
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

// func (r *repositories) SetTicketRedis(ctx context.Context, tickets []response.Ticket) error {
// 	// set data to redis
// 	val, err := json.Marshal(tickets)
// 	if err != nil {
// 		r.log.Error(ctx, "From Repositories: Failed to marshal data", err)
// 		return errors.BadRequest("Failed to marshal data")
// 	}

// 	if err := r.redisClient.Set(ctx, "tickets", val, 0).Err(); err != nil {
// 		r.log.Error(ctx, "From Repositories: Failed to set data to redis", err)
// 		return err
// 	}

// 	return nil
// }

// func (r *repositories) GetTicketRedis(ctx context.Context) ([]response.Ticket, error) {
// 	// get data from redis
// 	val, err := r.redisClient.Get(ctx, "tickets").Result()
// 	if err != nil {
// 		r.log.Error(ctx, "From Repositories: Failed to get data from redis", err)
// 		return nil, err
// 	}

// 	// parse response
// 	var tickets []response.Ticket
// 	if err := json.Unmarshal([]byte(val), &tickets); err != nil {
// 		r.log.Error(ctx, "From Repositories: Failed to unmarshal data", err)
// 		return nil, err
// 	}

// 	return tickets, nil
// }

func (r *repositories) ValidateToken(ctx context.Context, token string) (response.UserServiceValidate, error) {
	// http call to user service
	url := fmt.Sprintf("http://%s:%s/api/private/user/validate?token=%s", r.cfgUserService.Host, r.cfgUserService.Port, token)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response.UserServiceValidate{}, err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return response.UserServiceValidate{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to validate token, %v", resp.StatusCode))
		return response.UserServiceValidate{}, errors.BadRequest("Invalid token")
	}

	// parse response
	var respBase response.BaseResponse

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respBase); err != nil {
		return response.UserServiceValidate{
			IsValid: false,
			UserID:  0,
		}, err
	}

	jsonData, err := json.Marshal(respBase.Data)
	if err != nil {
		return response.UserServiceValidate{
			IsValid:   false,
			UserID:    0,
			EmailUser: "",
		}, err
	}

	var dataUser response.UserServiceValidate

	if err := json.Unmarshal(jsonData, &dataUser); err != nil {
		return response.UserServiceValidate{
			IsValid: false,
			UserID:  0,
		}, err
	}

	if !dataUser.IsValid {
		r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to validate token, %v", resp.StatusCode))
		return response.UserServiceValidate{
			IsValid: false,
			UserID:  0,
		}, errors.BadRequest("Invalid token")
	}

	// validate token
	return response.UserServiceValidate{
		IsValid:   dataUser.IsValid,
		UserID:    dataUser.UserID,
		EmailUser: dataUser.EmailUser,
	}, nil
}

// GetTicketOnline implements Repositories.
func (r *repositories) GetTicketOnline(ctx context.Context, regionName string) (response.OnlineTicket, error) {
	// http call to user service
	url := fmt.Sprintf("http://%s:%s/api/private/online-ticket?region_name=%s", r.cfgRecommendationService.Host, r.cfgRecommendationService.Port, regionName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response.OnlineTicket{}, err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return response.OnlineTicket{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get ticket online, %v", resp.StatusCode))
		return response.OnlineTicket{}, errors.BadRequest("Failed to get ticket online")
	}

	// parse response
	// var respData response.OnlineTicket

	// dec := json.NewDecoder(resp.Body)
	// if err := dec.Decode(&respData); err != nil {
	// 	return response.OnlineTicket{}, err
	// }

	// parse response
	var respBase response.BaseResponse

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respBase); err != nil {
		return response.OnlineTicket{}, err
	}

	respBase.Data = respBase.Data.(map[string]interface{})
	respData := response.OnlineTicket{
		IsSoldOut:      respBase.Data.(map[string]interface{})["is_sold_out"].(bool),
		IsFirstSoldOut: respBase.Data.(map[string]interface{})["is_first_sold_out"].(bool),
	}

	return respData, nil
}

// GetProfile implements Repositories.
func (r *repositories) GetProfile(ctx context.Context, userID int64) (response.Profile, error) {
	// http call to user service
	url := fmt.Sprintf("http://%s:%s/api/private/user/profile?user_id=%d", r.cfgUserService.Host, r.cfgUserService.Port, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response.Profile{}, err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return response.Profile{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		r.log.Ctx(ctx).Error(fmt.Sprintf("Failed to get profile, %v", resp.StatusCode))
		return response.Profile{}, errors.BadRequest("Failed to get profile")
	}

	// parse response
	var respBase response.BaseResponse

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respBase); err != nil {
		return response.Profile{}, err
	}

	respBase.Data = respBase.Data.(map[string]interface{})
	respData := response.Profile{
		ID:             int(respBase.Data.(map[string]interface{})["id"].(float64)),
		UserID:         int64(respBase.Data.(map[string]interface{})["user_id"].(float64)),
		FirstName:      respBase.Data.(map[string]interface{})["first_name"].(string),
		LastName:       respBase.Data.(map[string]interface{})["last_name"].(string),
		Address:        respBase.Data.(map[string]interface{})["address"].(string),
		District:       respBase.Data.(map[string]interface{})["district"].(string),
		City:           respBase.Data.(map[string]interface{})["city"].(string),
		State:          respBase.Data.(map[string]interface{})["state"].(string),
		Country:        respBase.Data.(map[string]interface{})["country"].(string),
		Region:         respBase.Data.(map[string]interface{})["region"].(string),
		Phone:          respBase.Data.(map[string]interface{})["phone"].(string),
		PersonalID:     respBase.Data.(map[string]interface{})["personal_id"].(string),
		TypePersonalID: respBase.Data.(map[string]interface{})["type_personal_id"].(string),
	}

	return respData, nil
}
