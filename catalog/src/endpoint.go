package src

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

//req-resp
type (
	CreateGolfRequest struct {
		id        string   `json:"id"`
		name      string   `json:"name"`
		location  string   `json:"location"`
		amenities []string `json:"amenities"`
	}
	CreateGolfResponse struct {
		ok string `json:"ok"`
	}

	GetGolfRequest struct {
		id string `json:"id"`
	}
	GetGolfResponse struct {
		id        string   `json:"id"`
		name      string   `json:"name"`
		location  string   `json:"location"`
		amenities []string `json:"amenities"`
	}
)

//Endpoints
func makeCreateGolfEndpoint(s GolfService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateGolfRequest)
		ok, err := s.CreateGolf(ctx, req.id, req.name, req.location, req.amenities)
		return CreateGolfResponse{ok: ok}, err
	}
}

func makeGetGolfEndpoint(s GolfService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetGolfRequest)
		golf, err := s.GetGolf(ctx, req.id)
		return GetGolfResponse{
			id:        golf.Id,
			name:      golf.Name,
			location:  golf.Location,
			amenities: golf.Amenities,
		}, err
	}
}
