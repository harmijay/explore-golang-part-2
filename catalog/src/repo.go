package src

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var RepoErr = errors.New("unable to handle repo request")

const GolfCollection = "golf"

type repo struct {
	db     *mongo.Database
	logger log.Logger
}

func NewRepo(db *mongo.Database, logger log.Logger) (Repository, error) {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "mongodb"),
	}, nil
}

func (repo *repo) CreateGolf(ctx context.Context, golf Golf) error {

	if golf.Id == "" || golf.Name == "" || golf.Location == "" || len(golf.Amenities) < 0 {
		return RepoErr
	}

	_, err := repo.db.Collection(GolfCollection).InsertOne(ctx, golf)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetGolf(ctx context.Context, id string) (Golf, error) {
	var golf Golf
	err := repo.db.Collection(GolfCollection).FindOne(ctx, bson.M{"id": id}).Decode(&golf)
	if err != nil {
		return Golf{}, RepoErr
	}

	return golf, nil
}
