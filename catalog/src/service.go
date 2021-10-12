package src

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var ErrInvalidArgument = errors.New("invalid argument")

type GolfService interface {
	CreateGolf(ctx context.Context, id string, name string, location string, amenities []string) (string, error)
	GetGolf(ctx context.Context, id string) (Golf, error)
}

type service struct {
	repository Repository
	logger     log.Logger
}

func NewService(rep Repository, logger log.Logger) GolfService {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s service) CreateGolf(ctx context.Context, id string, name string, location string, amenities []string) (string, error) {
	logger := log.With(s.logger, "method", "CreateGolf")

	golf := Golf{
		Id:        id,
		Name:      name,
		Location:  location,
		Amenities: amenities,
	}

	if err := s.repository.CreateGolf(ctx, golf); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("create golf", id)

	return "Success", nil
}

func (s service) GetGolf(ctx context.Context, id string) (Golf, error) {
	if id == "" {
		return Golf{}, ErrInvalidArgument
	}

	logger := log.With(s.logger, "method", "GetGolf")

	golf, err := s.repository.GetGolf(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return Golf{}, err
	}

	logger.Log("Get golf", id)

	return golf, nil
}
