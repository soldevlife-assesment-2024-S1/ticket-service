package messagestream

import (
	"fmt"
	"log"
	"net/url"
	"ticket-service/config"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

type ampq struct {
	cfg *config.MessageStreamConfig
}

func NewAmpq(cfg *config.MessageStreamConfig) MessageStream {
	return &ampq{
		cfg: cfg,
	}
}

func (m *ampq) NewSubscriber() (message.Subscriber, error) {
	username := url.QueryEscape(m.cfg.Username)
	password := url.QueryEscape(m.cfg.Password)
	host := m.cfg.Host
	port := m.cfg.Port

	var ampqURI string
	if m.cfg.SSL {
		ampqURI = fmt.Sprintf("amqps://%s:%s@%s:%s/", username, password, host, port)
	} else {
		ampqURI = fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
	}
	ampqConfig := amqp.NewDurableQueueConfig(ampqURI)

	subscriber, err := amqp.NewSubscriber(
		ampqConfig,
		watermill.NewStdLogger(stateLog, stateLog),
	)
	if err != nil {
		log.Fatal(err)
	}

	return subscriber, err
}

func (m *ampq) NewPublisher() (message.Publisher, error) {
	username := url.QueryEscape(m.cfg.Username)
	password := url.QueryEscape(m.cfg.Password)
	host := m.cfg.Host
	port := m.cfg.Port

	var ampqURI string
	if m.cfg.SSL {
		ampqURI = fmt.Sprintf("amqps://%s:%s@%s:%s/", username, password, host, port)
	} else {
		ampqURI = fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
	}
	ampqConfig := amqp.NewDurableQueueConfig(ampqURI)

	publisher, err := amqp.NewPublisher(
		ampqConfig,
		watermill.NewStdLogger(stateLog, stateLog),
	)
	if err != nil {
		log.Fatal(err)
	}

	tracePropagatingPublisherDecorator := wotelfloss.NewTracePropagatingPublisherDecorator(publisher)
	return wotel.NewNamedPublisherDecorator(m.cfg.ExchangeName, tracePropagatingPublisherDecorator), err
}

func ProcessMessages(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("Got message: %s", string(msg.Payload))
		msg.Ack()
	}
}
