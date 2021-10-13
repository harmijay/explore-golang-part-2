package src

import (
	"account/src/model"
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	CreateUser(ctx context.Context, username string, email string, password string) (string, error)
	GetUser(ctx context.Context, username string) (model.User, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewService(repo Repository, logger log.Logger) Service {
	return &service{
		repo: repo,
		logger: logger,
	}
}

func (s *service) CreateUser(ctx context.Context, username string,
							 email string, password string) (string, error)  {
	logger := log.With(s.logger, "method", "CreateUser")
	if username == "" || email == "" || password == "" {
		return "", ErrInvalidArgument
	}

	id := uuid.New().String()
	user := model.User{id, username, email, password}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("create user", username)
	return "Success", nil
}

func (s *service) GetUser(ctx context.Context, username string) (model.User, error) {
	logger := log.With(s.logger, "method", "GetUser")

	if username == "" {
		return model.User{}, ErrInvalidArgument
	}

	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		level.Error(logger).Log("err", err)
	}

	logger.Log("get user", username)
	return user, nil
}
