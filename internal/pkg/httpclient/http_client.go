package httpclient

import (
	"net/http"
	"ticket-service/config"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	circuit "github.com/rubyist/circuitbreaker"
)

// Init initializes the circuit breaker based on the configuration and breaker type: consecutive, error_rate, threshold
func InitCircuitBreaker(cfg *config.HttpClientConfig, breakerType string) (cb *circuit.Breaker) {
	switch breakerType {
	case "consecutive":
		cb = circuit.NewConsecutiveBreaker(
			int64(cfg.ConsecutiveFailures),
		)
	case "error_rate":
		cb = circuit.NewRateBreaker(
			cfg.ErrorRate, 100,
		)
	default:
		if cfg.Threshold == 0 {
			cfg.Threshold = 10
		}
		cb = circuit.NewThresholdBreaker(
			int64(cfg.Threshold),
		)
	}
	return cb
}

func InitHttpClient(cfg *config.HttpClientConfig, cb *circuit.Breaker) *circuit.HTTPClient {
	timeout := time.Duration(cfg.Timeout) * time.Second
	client := circuit.NewHTTPClientWithBreaker(
		cb,
		timeout,
		nil,
	)
	client.Client.Transport = otelhttp.NewTransport(http.DefaultTransport)
	return client
}
