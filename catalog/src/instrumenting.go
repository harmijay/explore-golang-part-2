package src

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	GolfService
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s GolfService) GolfService {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		GolfService:    s,
	}
}

func (s *instrumentingService) CreateGolf(ctx context.Context, id string, name string, location string, amenities []string) (msg string, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create_golf").Add(1)
		s.requestLatency.With("method", "create_golf").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.GolfService.CreateGolf(ctx, id, name, location, amenities)
}

func (s *instrumentingService) GetGolf(ctx context.Context, id string) (g Golf, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get_golf").Add(1)
		s.requestLatency.With("method", "get_golf").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.GolfService.GetGolf(ctx, id)
}
