package src

import (
	"account/src/model"
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   	metrics.Counter
	requestLatency 	metrics.Histogram
	service			Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, service Service) Service {
	return &instrumentingService{
		requestCount: counter,
		requestLatency: latency,
		service: service,
	}
}

func (s *instrumentingService) CreateUser(ctx context.Context, username string,
										  email string, password string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create user").Add(1)
		s.requestLatency.With("method", "create user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.service.CreateUser(ctx, username, email, password)
}

func (s *instrumentingService) GetUser(ctx context.Context, username string) (model.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.service.GetUser(ctx, username)
}
