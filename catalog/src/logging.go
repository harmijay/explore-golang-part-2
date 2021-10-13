package src

import (
	"context"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	GolfService
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s GolfService) GolfService {
	return &loggingService{logger, s}
}

func (s *loggingService) CreateGolf(ctx context.Context, id string, name string, location string, amenities []string) (msg string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "create_golf",
			"id", id,
			"name", name,
			"location", location,
			"amenities", "["+strings.Join(amenities, ", ")+"]",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.GolfService.CreateGolf(ctx, id, name, location, amenities)
}

func (s *loggingService) GetGolf(ctx context.Context, id string) (g Golf, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get_golf",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.GolfService.GetGolf(ctx, id)
}
