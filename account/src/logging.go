package src

import (
	"account/src/model"
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	service Service
}

func NewLoggingService(logger log.Logger, service Service) Service {
	return &loggingService{
		logger: logger,
		service: service,
	}
}

func (s *loggingService) CreateUser(ctx context.Context, username string,
									email string, password string) (string, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "create user",
			"username", username,
			"email", email,
			"password", password,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.service.CreateUser(ctx, username, email, password)
}

func (s *loggingService) GetUser(ctx context.Context, username string) (model.User, error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get user",
			"username", username,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.service.GetUser(ctx, username)
}
