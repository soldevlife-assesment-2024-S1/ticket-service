package messagestream

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"ticket-service/internal/module/ticket/models/request"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

var (
	stateLog, _ = strconv.ParseBool(os.Getenv("PRODUCTION"))
)

type MessageStream interface {
	NewSubscriber() (message.Subscriber, error)
	NewPublisher() (message.Publisher, error)
}

func NewRouter(pub message.Publisher, poisonTopic string, handlerTopicName string, subscribeTopic string, subs message.Subscriber, handlerFunc func(msg *message.Message) error) (*message.Router, error) {
	logger := watermill.NewStdLogger(stateLog, stateLog)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	router.AddPlugin(plugin.SignalsHandler)

	poisonMiddleware, err := middleware.PoisonQueue(
		pub,
		poisonTopic,
	)

	if err != nil {
		return nil, err
	}

	router.AddMiddleware(
		middleware.CorrelationID,
		poisonMiddleware,

		middleware.Retry{
			MaxRetries:      5,
			InitialInterval: 500,
			Logger:          logger,
		}.Middleware,

		middleware.CorrelationID,
		middleware.Recoverer,
		wotel.Trace(),
	)

	router.AddNoPublisherHandler(
		handlerTopicName,
		subscribeTopic,
		subs,
		handlerFunc,
	)

	router.AddPlugin(plugin.SignalsHandler)

	return router, err
}

func PoisonedQueue(err error, p message.Publisher, msg *message.Message, topicTarget string, log *otelzap.Logger) {
	// publish to poison queue
	reqPoisoned := request.PoisonedQueue{
		TopicTarget: topicTarget,
		ErrorMsg:    err.Error(),
		Payload:     msg.Payload,
	}

	jsonPayload, _ := json.Marshal(reqPoisoned)

	err = p.Publish("poisoned_queue", message.NewMessage(watermill.NewUUID(), jsonPayload))
	if err != nil {
		log.Ctx(msg.Context()).Error(fmt.Sprintf("Failed to publish to poison queue: %v", err))
	}

}
