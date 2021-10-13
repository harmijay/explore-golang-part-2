package src

import (
	"account/src/model"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type CreateUserRequest struct {
	Username	string	`json:"username" binding:"required"`
	Email		string	`json:"email" binding:"required"`
	Password	string	`json:"password" binding:"required"`
}

type CreateUserResponse struct {
	Ok		string	`json:"ok,omitempty"`
	Err		error	`json:"err,omitempty"`
}

func (r CreateUserResponse) error() error {
	return r.Err
}

func MakeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		ok, err := s.CreateUser(ctx, req.Username, req.Email, req.Password)
		return CreateUserResponse{Ok: ok, Err: err}, nil
	}
}

type GetUserRequest struct {
	Username	string	`json:"username" binding:"required"`
}

type GetUserResponse struct {
	User	model.User	`json:"user,omitempty"`
	Err		error		`json:"err,omitempty"`
}

func (r GetUserResponse) error() error {
	return r.Err
}

func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := s.GetUser(ctx, req.Username)
		return GetUserResponse{User: user, Err: err}, nil
	}
}
