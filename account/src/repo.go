package src

import (
	"account/src/model"
	"context"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, username string) (model.User, error)
}

type repo struct {
	db		*mongo.Database
	logger	log.Logger
}

func NewRepo(db *mongo.Database, logger log.Logger) Repository {
	return &repo{
		db: db,
		logger: logger,
	}
}

func (repo *repo) CreateUser(ctx context.Context, user model.User) error {
	_, err := repo.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetUser(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := repo.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
