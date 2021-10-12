package src

import "context"

type Golf struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Location  string   `json:"location"`
	Amenities []string `json:"amenities"`
}

type Repository interface {
	CreateGolf(ctx context.Context, golf Golf) error
	GetGolf(ctx context.Context, id string) (Golf, error)
}
